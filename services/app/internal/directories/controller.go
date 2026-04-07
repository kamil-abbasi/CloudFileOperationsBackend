package directories

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/patrickmn/go-cache"
)

type DirectoriesController struct {
	config             *config.Config
	cache              *cache.Cache
	directoriesService *DirectoriesService
}

func NewController(config *config.Config, cache *cache.Cache, directoriesService *DirectoriesService) *DirectoriesController {
	return &DirectoriesController{
		cache:              cache,
		config:             config,
		directoriesService: directoriesService,
	}
}

func (c *DirectoriesController) ListItems(ctx *gin.Context) {
	location := ctx.Query("location")

	if location == "" {
		location = "/"
	}

	location = filepath.Clean(location)

	items, err := c.directoriesService.ListItems(location)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	// prevents from returning null when array is empty
	if len(items) == 0 {
		items = []dtos.DirectoryItemResponseDto{}
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *DirectoriesController) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")

	dir, err := c.directoriesService.FindOne(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, dir)
}

func (c *DirectoriesController) Download(ctx *gin.Context) {
	ctx.Header("Content-Disposition", "attachment; filename=download.zip")
	ctx.Header("Content-Type", "application/zip")

	id := ctx.Param("id")

	err := c.directoriesService.Download(id, ctx.Writer)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

}

func (c *DirectoriesController) Create(ctx *gin.Context) {
	data, _ := ctx.Get("validatedBody")
	createDto := data.(dtos.CreateDirectoryDto)
	userId := ctx.GetString("userId")

	directory, err := c.directoriesService.Create(userId, createDto)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, directory)
}

func (c *DirectoriesController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	data, _ := ctx.Get("validatedBody")
	dto := data.(dtos.UpdateDirectoryDto)

	directory, err := c.directoriesService.Update(id, dto)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, directory)
}

func (c *DirectoriesController) Remove(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.directoriesService.Remove(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.Status(http.StatusOK)
}
