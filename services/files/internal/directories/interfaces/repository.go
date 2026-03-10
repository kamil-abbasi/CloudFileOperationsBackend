package interfaces

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
)

type IDirectoriesRepository interface {
	Find() ([]entities.Directory, error)
	FindOne(id string) (entities.Directory, bool, error)
	Create(dto dtos.CreateDirectoryDto) (entities.Directory, error)
	Remove(id string) (bool, error)
}
