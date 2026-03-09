package files

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type DirectoriesController struct {
	config       *config.Config
	filesService *files.FilesService
}

func NewDirectoriesController(config *config.Config) DirectoriesController {
	return DirectoriesController{
		config:       config,
		filesService: files.NewService(config),
	}
}

// POST /v1/directories?path=/some/folder
func (c *DirectoriesController) Create(ctx *gin.Context) {
	path := ctx.Query("path")
	userId := "user-dev"

	if !isValidPath(path) {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid or missing path",
		})

		return
	}

	cleanPath := cleanLocation(path)

	err := c.filesService.CreateDirectory(userId, cleanPath)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Created directory %v", cleanPath),
		"path":    cleanPath,
	})
}

// DELETE /v1/directories?path=/some/folder
func (c *DirectoriesController) Remove(ctx *gin.Context) {
	path := ctx.Query("path")
	userId := "user-dev"

	if !isValidPath(path) {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid or missing path",
		})

		return
	}

	cleanPath := cleanLocation(path)

	removedCount, err := c.filesService.RemoveDirectory(userId, cleanPath)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	if removedCount == 0 {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "Directory not found or is empty",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      fmt.Sprintf("Removed directory %v", cleanPath),
		"removedFiles": removedCount,
	})
}

// GET /v1/directories/download?path=/some/folder
func (c *DirectoriesController) Download(ctx *gin.Context) {
	path := ctx.Query("path")
	userId := "user-dev"

	if !isValidPath(path) {
		ctx.JSON(http.StatusBadRequest, &shared.HttpError{
			Code:    http.StatusBadRequest,
			Message: "Invalid or missing path",
		})

		return
	}

	cleanPath := cleanLocation(path)

	files, err := c.filesService.FindByLocation(userId, cleanPath)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	if len(files) == 0 {
		ctx.JSON(http.StatusNotFound, &shared.HttpError{
			Code:    http.StatusNotFound,
			Message: "Directory not found or is empty",
		})

		return
	}

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v.zip\"", filepath.Base(cleanPath)))

	err = c.filesService.DownloadDirectory(userId, cleanPath, ctx.Writer)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}
}

func isValidPath(path string) bool {
	if path == "" {
		return false
	}

	re, err := regexp.Compile(`^\/?([\w-]+\/?)+$`)

	if err != nil {
		return false
	}

	return len(re.FindStringIndex(path)) > 0
}

func cleanLocation(path string) string {
	cleaned := filepath.Clean(path)
	cleaned = strings.TrimPrefix(cleaned, "/")
	cleaned = strings.TrimPrefix(cleaned, "\\")

	return cleaned
}
