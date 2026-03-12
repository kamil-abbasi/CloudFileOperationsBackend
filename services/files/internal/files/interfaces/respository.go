package interfaces

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IFilesRepository interface {
	Find(dto dtos.FindFilesDto) ([]entities.File, error)
	FindOne(id string) (entities.File, bool, error)
	FindByNameAndDirectoryId(name string, directoryId string) (entities.File, bool, error)
	Save(file entities.File) error
	Remove(id string) (bool, error)
}
