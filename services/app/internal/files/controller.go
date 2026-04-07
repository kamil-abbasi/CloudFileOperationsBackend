package files

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
)

type FilesController struct {
	config       *config.Config
	cache        *cache.Cache
	filesService *FilesService
}

func NewController(config *config.Config, cache *cache.Cache, filesService *FilesService) *FilesController {
	return &FilesController{
		cache:        cache,
		config:       config,
		filesService: filesService,
	}
}

// POST /v1/files
func (c *FilesController) Upload(ctx *gin.Context) {
	data, _ := ctx.Get("validatedBody")
	dto := data.(dtos.UploadFileDto)

	userId := ctx.GetString("userId")

	requestId := dto.IdempotencyKey + userId
	response, found := c.cache.Get(requestId)

	if found {
		ctx.JSON(http.StatusCreated, response)

		return
	}

	reader, err := dto.File.Open()

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	defer reader.Close()

	file, err := c.filesService.Create(userId, dtos.CreateFileDto{
		Name:        dto.File.Filename,
		DirectoryId: dto.DirectoryId,
	}, reader)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	c.cache.Set(requestId, file, 15*time.Minute)
	ctx.JSON(http.StatusCreated, file)
}

// GET /v1/files/:id/download
func (c *FilesController) Download(ctx *gin.Context) {
	id := ctx.Param("id")

	file, err := c.filesService.FindOne(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	src, err := c.filesService.Download(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	defer src.Close()

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%v\"", file.Name),
	}
	contentType := mime.TypeByExtension(filepath.Ext(file.Name))

	// sending file
	ctx.DataFromReader(http.StatusCreated, int64(file.Size), contentType, src, extraHeaders)
}

// GET /v1/files/:id
func (c *FilesController) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")

	file, err := c.filesService.FindOne(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// PATCH /v1/files/:id
func (c *FilesController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	data, _ := ctx.Get("validatedBody")
	updateDto := data.(dtos.UpdateFileDto)

	file, err := c.filesService.Update(id, updateDto)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DELETE /v1/files/:id
func (c *FilesController) Remove(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.filesService.Remove(id)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.Status(http.StatusOK)
}
