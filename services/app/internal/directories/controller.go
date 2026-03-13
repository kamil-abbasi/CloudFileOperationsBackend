package directories

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
	"github.com/patrickmn/go-cache"
)

type DirectoriesController struct {
	config             *config.Config
	cache              *cache.Cache
	directoriesService *DirectoriesService
}

func NewController(config *config.Config, cache *cache.Cache, service *DirectoriesService) DirectoriesController {
	return DirectoriesController{
		cache:              cache,
		config:             config,
		directoriesService: service,
	}
}

func (c *DirectoriesController) FindOne(ctx *gin.Context) {}

func (c *DirectoriesController) Download(ctx *gin.Context) {
	ctx.Header("Content-Disposition", "attachment; filename=download.zip")
	ctx.Header("Content-Type", "application/zip")

	id := ctx.Param("id")

	found, err := c.directoriesService.Download(id, ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	if !found {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "Directory not found",
		})

		return
	}

}

func (c *DirectoriesController) Create(ctx *gin.Context) {
	var fields struct {
		Name     string
		ParentId string
	}

	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadGateway,
			Message: "Invalid request body",
		})

		return
	}

	createDto := dtos.CreateDirectoryDto{
		UserId:   "user-dev",
		ParentId: fields.ParentId,
		Name:     fields.Name,
	}

	directory, err := c.directoriesService.Create(createDto)

	if err != nil {
		// TODO: handle errs
		_, ok := err.(*shared.DirectoryAlreadyExistsError)

		if ok {
			ctx.JSON(http.StatusConflict, &shared.HttpError{
				Code:    http.StatusConflict,
				Message: "Directory already exists",
			})

			return
		}

		_, ok = err.(*shared.DirectoryNotFoundError)

		if ok {
			ctx.JSON(http.StatusNotFound, &shared.HttpError{
				Code:    http.StatusNotFound,
				Message: "Parent directory not found",
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal error",
		})

		return
	}

	ctx.JSON(http.StatusOK, directory)
}

func (c *DirectoriesController) Remove(ctx *gin.Context) {
	id := ctx.Param("id")

	removed, err := c.directoriesService.Remove(id)

	if err != nil {
		log.Print(err)

		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	if !removed {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "Directory not found",
		})

		return
	}

	ctx.Status(http.StatusOK)
}
