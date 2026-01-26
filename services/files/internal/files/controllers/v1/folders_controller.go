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

type FoldersController struct {
	config       *config.Config
	filesService *files.FilesService
}

func NewFoldersController(config *config.Config) FoldersController {
	return FoldersController{
		config:       config,
		filesService: files.NewService(config),
	}
}

// POST /v1/folders
func (c *FoldersController) Create(ctx *gin.Context) {
	var body struct {
		Name     string
		Location string
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	userId := "user-dev"

	id, err := uuid.NewUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	folder, err := c.filesService.CreateFolder(dtos.FolderCreateDto{
		Id:       id.String(),
		UserId:   userId,
		Name:     body.Name,
		Location: body.Location,
	})

	if err != nil {
		_, ok := err.(*errs.FolderAlreadyExistsError)
		if ok {
			ctx.JSON(http.StatusConflict, &shared.HttpError{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("Folder '%s' already exists", body.Name),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, folder)
}
