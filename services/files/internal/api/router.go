package api

import (
	"github.com/gin-gonic/gin"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
)

func NewRouter(filesController *files.FilesController) *gin.Engine {
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
	}

	return r
}
