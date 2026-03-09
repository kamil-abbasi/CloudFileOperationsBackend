package interfaces

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IFilesRepository interface {
	FindOne(id string) (entities.File, bool, error)
	FindByLocation(userId string, location string) ([]entities.File, error)
	Update(dto dtos.FileUpdateDto) (bool, error)
	Remove(id string) (bool, error)
	RemoveByLocation(userId string, location string) (int64, error)
	Create(dto dtos.FileCreateDto) (entities.File, error)
}
