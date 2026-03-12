package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/api"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
	"github.com/patrickmn/go-cache"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Failed to start app, details: %v", err)
	}
}

func Run() error {
	cfg, err := config.Load()

	if err != nil {
		return err
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

	defer postgres.Close()

	cache := cache.New(15*time.Minute, 10*time.Minute)

	storageAdapter := storage.NewFileSystemStorageAdapter(cfg)

	filesRepository := files.NewPostgresRepository(cfg, postgres)
	directoriesRepository := directories.NewPostgresRepository(cfg, postgres)

	filesService := files.NewService(filesRepository, directoriesRepository, storageAdapter)
	directoriesService := directories.NewService(cfg, directoriesRepository, filesRepository)

	filesController := files.NewController(cfg, cache, &filesService)
	directoriesController := directories.NewController(cfg, cache, &directoriesService)

	r := api.NewRouter(&filesController, &directoriesController)
	r.SetTrustedProxies([]string{"reverse-proxy"})
	r.Run()

	return nil
}
