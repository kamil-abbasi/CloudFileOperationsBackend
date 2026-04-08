package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/api"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/cache"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/usage"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Failed to start app, details: %v", err)
	}
}

func Run() error {
	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("Failed to load app configuration. Details: %v", err)
	}

	// creating root path for user files
	err = os.MkdirAll(cfg.RootPath, 0755)

	if err != nil {
		return fmt.Errorf("Failed to create directory for user files. Details: %v", err)
	}

	postgres, err := database.NewPostgres(database.PostgresConfig{
		Db:       cfg.PgDb,
		User:     cfg.PgUser,
		Password: cfg.PgPassword,
		Host:     cfg.PgHost,
	})

	if err != nil {
		return err
	}

	defer postgres.Conn.Close()

	// general
	cache := cache.NewInMemoryAdapter()
	storage := storage.NewFileSystemAdapter(cfg)

	// repos
	filesRepository := files.NewPostgresRepository(postgres)
	directoriesRepository := directories.NewPostgresRepository(postgres)
	usersRepository := auth.NewPostgresRepository(postgres)

	// services
	filesService := files.NewService(filesRepository, directoriesRepository, storage)
	directoriesService := directories.NewService(cfg, directoriesRepository, filesRepository, storage)
	usageService := usage.NewService(filesService)
	authService := auth.NewService(usersRepository, cfg)

	// controllers
	filesController := files.NewController(cfg, cache, filesService)
	directoriesController := directories.NewController(cfg, cache, directoriesService)
	usageController := usage.NewController(usageService)
	authController := auth.NewController(authService)

	r := api.NewRouter(api.RouterParams{
		FilesController:       filesController,
		DirectoriesController: directoriesController,
		UsageController:       usageController,
		AuthController:        authController,
		Cfg:                   cfg,
	})
	err = r.SetTrustedProxies([]string{"reverse-proxy"})
	err = r.Run()

	if err != nil {
		return err
	}

	return nil
}
