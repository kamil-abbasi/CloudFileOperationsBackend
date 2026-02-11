package v1

import (
	"github.com/gin-gonic/gin"

	files "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
)

func RegisterRoutes(rg *gin.RouterGroup, filesController *files.FilesController) {
	filesRouter := rg.Group("/files")
	{
		filesRouter.GET("/:id", filesController.FindOne)

		filesRouter.GET("/:id/download", filesController.Download)

		filesRouter.POST("", filesController.Upload)

		filesRouter.PATCH("/:id", filesController.Update)

		filesRouter.DELETE("/:id", filesController.Remove)
	}

}
