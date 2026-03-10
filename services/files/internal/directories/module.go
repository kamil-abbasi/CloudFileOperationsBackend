package directories

import (
	"database/sql"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	fileInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/patrickmn/go-cache"
)

type ModuleDeps struct {
	Db              *sql.DB
	Config          *config.Config
	Cache           *cache.Cache
	FilesRepository fileInterfaces.IFilesRepository
}

type ModuleExports struct {
	Repository *interfaces.IDirectoriesRepository
	Service    *DirectoriesService
	Controller *DirectoriesController
}

func NewModule(deps *ModuleDeps) *ModuleExports {
	repository := NewPostgresRepository(deps.Config, deps.Db)
	service := NewService(deps.Config, repository, deps.FilesRepository)
	controller := NewController(deps.Config, deps.Cache, &service)

	return &ModuleExports{
		Repository: &repository,
		Service:    &service,
		Controller: &controller,
	}
}
