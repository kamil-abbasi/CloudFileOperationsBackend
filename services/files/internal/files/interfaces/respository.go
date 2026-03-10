package interfaces

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IFilesRepository interface {
	Find(dto dtos.FindFilesDto) ([]entities.File, error)
	FindOne(id string) (entities.File, bool, error)
	Create(dto dtos.CreateFileDto) (entities.File, error)
	Update(dto dtos.UpdateFileDto) (bool, error)
	Remove(id string) (bool, error)
}
