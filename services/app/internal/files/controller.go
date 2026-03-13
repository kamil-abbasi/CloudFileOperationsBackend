package files

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FilesController struct {
	config       *config.Config
	cache        *cache.Cache
	filesService *FilesService
}

func NewController(config *config.Config, cache *cache.Cache, service *FilesService) FilesController {
	return FilesController{
		cache:        cache,
		config:       config,
		filesService: service,
	}
}

// POST /v1/files
func (c *FilesController) Upload(ctx *gin.Context) {
	rawFile, err := ctx.FormFile("file")
	directoryId := ctx.PostForm("directoryId")
	userId := "user-dev"

	idempotencyKey := ctx.PostForm("idempotencyKey")
	requestId := idempotencyKey + userId

	if idempotencyKey == "" {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Missing idempotency key",
		})

		return
	}

	response, found := c.cache.Get(requestId)

	if found {
		ctx.JSON(http.StatusCreated, response)

		return
	}

	reader, err := rawFile.Open()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	defer reader.Close()

	file, err := c.filesService.Create(dtos.CreateFileDto{
		Name:        rawFile.Filename,
		UserId:      userId,
		DirectoryId: directoryId,
	}, reader)

	if err != nil {

		_, ok := err.(*shared.FileAlreadyExistsError)

		if ok {
			ctx.JSON(http.StatusConflict, &shared.HttpError{
				Code:    http.StatusConflict,
				Message: "File already exists",
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	c.cache.Set(requestId, file, 15*time.Minute)
	ctx.JSON(http.StatusCreated, file)
}

// GET /v1/files/:id/download
func (c *FilesController) Download(ctx *gin.Context) {
	id := ctx.Param("id")

	file, found, err := c.filesService.FindOne(id)

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
			Message: "File not found",
		})

		return
	}

	src, found, err := c.filesService.Download(id)

	if err != nil || !found || src == nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

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

	file, found, err := c.filesService.FindOne(id)

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
			Message: "File not found",
		})

		return
	}

	ctx.JSON(http.StatusOK, file)
}

// PATCH /v1/files/:id
func (c *FilesController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var fields struct {
		Name        string
		DirectoryId string
	}

	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})

		return
	}

	updateDto := dtos.UpdateFileDto{}
	updateDto.Where.Id = id
	updateDto.Fields.Name = fields.Name
	updateDto.Fields.DirectoryId = fields.DirectoryId

	file, updated, err := c.filesService.Update(updateDto)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	if !updated {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "File not found",
		})

		return
	}

	ctx.JSON(http.StatusOK, file)
}

// DELETE /v1/files/:id
func (c *FilesController) Remove(ctx *gin.Context) {
	id := ctx.Param("id")

	removed, err := c.filesService.Remove(id)

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
			Message: "File not found",
		})

		return
	}

	ctx.Status(http.StatusOK)
}
