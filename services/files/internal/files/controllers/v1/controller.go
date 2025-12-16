package files

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	errs "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/errors"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FilesController struct {
	config       *config.Config
	filesService *files.FilesService
}

func NewController(config *config.Config) FilesController {
	return FilesController{
		config:       config,
		filesService: files.NewService(config),
	}
}

// POST /v1/files
func (c *FilesController) Upload(ctx *gin.Context) {
	rawFile, err := ctx.FormFile("file")
	location := ctx.PostForm("location")
	userName := "user-dev"

	userRoot := fmt.Sprintf("%v/%v", c.config.RootPath, userName)
	fullPath := fmt.Sprintf("%v/%v/%v", userRoot, location, rawFile.Filename)
	userPath := fmt.Sprintf("/%v/%v", userName, location)

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
	})

	if err != nil {

		_, ok := err.(*errs.FileAlreadyExistsError)

		if ok {
			ctx.JSON(http.StatusConflict, &shared.HttpError{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("File %v already exists", userPath),
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

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

	// setting filename and forcing clients to download
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", file.Filename))

	// sending file
	filePath := c.config.RootPath + file.Location + "/" + file.Filename

	fmt.Println(filePath)

	ctx.File(filePath)
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
	}

	if !found {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "File not found",
		})
	}

	ctx.JSON(http.StatusOK, file)
}

// PATCH /v1/files/:id
func (c *FilesController) Update(ctx *gin.Context) {}

// DELETE /v1/files/:id
func (c *FilesController) Remove(ctx *gin.Context) {}
