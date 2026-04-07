package directories

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	filesEntities "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IDirectoriesRepository interface {
	ListItems(location string) ([]entities.DirectoryItem, error)
	Find() ([]entities.Directory, error)
	FindOne(id string) (entities.Directory, bool, error)
	FindByNameAndParentId(name string, parentId string) (entities.Directory, bool, error)
	Save(directory entities.Directory) error
	Remove(id string) (bool, error)
}

type IFilesRepository interface {
	FindByLocation(location string) ([]filesEntities.File, error)
}
