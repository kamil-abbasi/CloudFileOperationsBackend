package files

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type FileService struct {
	repository interfaces.IFilesRepository
}

func NewService() *FileService {

	repository := NewSQLiteRepository()

	return &FileService{
		repository: repository,
	}
}

func (service *FileService) FindOne(id string) (entities.File, bool, error) {
	return service.repository.FindOne(id)
}

func (service *FileService) Create(dto dtos.FileCreateDto) (entities.File, error) {
	return service.repository.Create(dto)
}

func (service *FileService) Update(dto dtos.FileUpdateDto) (entities.File, bool, error) {
	return service.repository.Update(dto)
}

func (service *FileService) Remove(id string) (entities.File, bool, error) {
	return service.repository.Remove(id)
}
