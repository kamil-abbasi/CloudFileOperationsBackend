package main

import (
	"log"
	"os"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/api"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
)

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Failed to start app, details: %v", err)
	}
}

func Run() error {
	cfg := config.New()
	cfg.Load()

	os.MkdirAll(cfg.RootPath, 0755)

	db := database.New(cfg)
	db.Init()

	r := api.NewRouter(cfg)
	r.SetTrustedProxies([]string{"reverse-proxy"})
	r.Run()

	return nil
}
