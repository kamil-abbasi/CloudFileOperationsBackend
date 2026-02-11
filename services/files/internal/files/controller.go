package files

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	errs "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/errors"
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
	var pathRegex, err = regexp.Compile(`^(\/[\w-]+)*([\w-])*$`)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	rawFile, err := ctx.FormFile("file")
	location := ctx.PostForm("location")
	idempotencyKey := ctx.PostForm("idempotencyKey")
	userId := "user-dev"
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

	if len(pathRegex.FindStringIndex(location)) == 0 {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid location",
		})

		return
	}

	fullPath := filepath.Clean(filepath.Join(
		c.config.RootPath,
		userId,
		location,
		rawFile.Filename,
	))
	userPath := filepath.Clean(location)
	userFullPath := filepath.Clean(filepath.Join(location, rawFile.Filename))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	id, err := uuid.NewUUID()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	file, err := c.filesService.Create(dtos.FileCreateDto{
		Id:       id.String(),
		Filename: rawFile.Filename,
		Location: userPath,
		Size:     uint64(rawFile.Size),
		UserId:   userId,
	})

	if err != nil {

		_, ok := err.(*errs.FileAlreadyExistsError)

		if ok {
			ctx.JSON(http.StatusConflict, &shared.HttpError{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("File %v already exists", userFullPath),
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

	ctx.SaveUploadedFile(rawFile, fullPath)
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

	// setting filename and informing clients to download
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", file.Filename))

	fullPath := filepath.Clean(filepath.Join(
		c.config.RootPath,
		file.UserId,
		file.Location,
		file.Filename,
	))

	// sending file
	ctx.File(fullPath)
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
		Filename string
		Location string
	}

	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})

		return
	}

	updateDto := dtos.FileUpdateDto{}
	updateDto.Where.Id = id
	updateDto.Fields = fields

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

	file, removed, err := c.filesService.Remove(id)

	if err != nil {
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

	ctx.JSON(http.StatusOK, file)
}
