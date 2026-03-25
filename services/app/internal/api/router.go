package api

import (
	"github.com/gin-gonic/gin"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
)

func NewRouter(filesController *files.FilesController, directoriesController *directories.DirectoriesController) *gin.Engine {
	r := gin.Default()

	r.GET("/healthcheck")

	v1 := r.Group("/v1")
	{
		filesRouter := v1.Group("/files")
		{
			filesRouter.GET("/:id", filesController.FindOne)

			filesRouter.GET("/:id/download", filesController.Download)

			filesRouter.POST("", filesController.Upload)

			filesRouter.PATCH("/:id", filesController.Update)

			filesRouter.DELETE("/:id", filesController.Remove)
		}

		directoriesRouter := v1.Group("/directories")
		{
			directoriesRouter.GET("/:id/download", directoriesController.Download)

			directoriesRouter.GET("/:id", directoriesController.FindOne)

			directoriesRouter.PATCH("/:id", directoriesController.Update)

			directoriesRouter.POST("", directoriesController.Create)

			directoriesRouter.DELETE("/:id", directoriesController.Remove)

			directoriesRouter.GET("/items", directoriesController.ListItems)
		}
	}

	return r
}
