package files

import (
	"database/sql"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/patrickmn/go-cache"
)

type ModuleDeps struct {
	Db     *sql.DB
	Config *config.Config
	Cache  *cache.Cache
}

type ModuleExports struct {
	Repository *interfaces.IFilesRepository
	Service    *FilesService
	Controller *FilesController
}

func NewModule(deps *ModuleDeps) *ModuleExports {
	repository := NewPostgresRepository(deps.Config, deps.Db)
	service := NewService(deps.Config, repository)
	controller := NewController(deps.Config, deps.Cache, &service)

	return &ModuleExports{
		Repository: &repository,
		Service:    &service,
		Controller: &controller,
	}
}
