package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	files "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/controllers/v1"
)

func RegisterRoutes(rg *gin.RouterGroup, config *config.Config) {
	filesRouter := rg.Group("/files")
	{
		controller := files.NewController(config)

		filesRouter.GET("/:id", controller.FindOne)

		filesRouter.GET("/:id/download", controller.Download)

		filesRouter.POST("", controller.Upload)

		filesRouter.PATCH("/:id", controller.Update)

		filesRouter.DELETE("/:id", controller.Remove)
	}

	foldersRouter := rg.Group("/folders")
	{
		foldersController := files.NewFoldersController(config)

		foldersRouter.POST("", foldersController.Create)
	}

}
