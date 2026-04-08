package api

import (
	"github.com/gin-gonic/gin"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/api/middleware"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth"
	authDtos "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	dirDtos "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	fileDtos "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/usage"
)

type RouterParams struct {
	FilesController       *files.FilesController
	DirectoriesController *directories.DirectoriesController
	UsageController       *usage.UsageController
	AuthController        *auth.AuthController
	Cfg                   *config.Config
}

func NewRouter(params RouterParams) *gin.Engine {
	r := gin.Default()

	filesController := params.FilesController
	directoriesController := params.DirectoriesController
	usageController := params.UsageController
	authController := params.AuthController
	cfg := params.Cfg

	r.Use(middleware.Error())

	r.GET("/healthcheck")

	v1 := r.Group("/v1")
	{
		filesRouter := v1.Group("/files")
		{
			filesRouter.Use(middleware.Auth(cfg))

			filesRouter.GET("/:id", filesController.FindOne)

			filesRouter.GET("/:id/download", filesController.Download)

			filesRouter.POST("", middleware.Validation[fileDtos.UploadFileDto](), filesController.Upload)

			filesRouter.PATCH("/:id", middleware.Validation[fileDtos.UpdateFileDto](), filesController.Update)

			filesRouter.DELETE("/:id", filesController.Remove)
		}

		directoriesRouter := v1.Group("/directories")
		{
			directoriesRouter.Use(middleware.Auth(cfg))

			directoriesRouter.GET("/:id/download", directoriesController.Download)

			directoriesRouter.GET("/:id", directoriesController.FindOne)

			directoriesRouter.POST("", middleware.Validation[dirDtos.CreateDirectoryDto](), directoriesController.Create)

			directoriesRouter.PATCH("/:id", middleware.Validation[dirDtos.UpdateDirectoryDto](), directoriesController.Update)

			directoriesRouter.DELETE("/:id", directoriesController.Remove)

			directoriesRouter.GET("/items", directoriesController.ListItems)
		}

		usageRouter := v1.Group("/usage")
		{
			usageRouter.Use(middleware.Auth(cfg))

			usageRouter.GET("", usageController.CalculateForLocation)
		}

		authRouter := v1.Group("/auth")
		{
			authRouter.POST("/register", middleware.Validation[authDtos.RegisterUserDto](), authController.Register)

			authRouter.POST("/login", middleware.Validation[authDtos.LoginUserDto](), authController.Login)
		}
	}

	return r
}
