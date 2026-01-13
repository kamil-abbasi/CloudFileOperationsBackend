package interfaces

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IFilesRepository interface {
	FindOne(id string) (entities.File, bool, error)
	Update(dto dtos.FileUpdateDto) (bool, error)
	Remove(id string) (bool, error)
	Create(dto dtos.FileCreateDto) (entities.File, error)
	CreateFolder(dto dtos.FolderCreateDto) (entities.Folder, error)
}
