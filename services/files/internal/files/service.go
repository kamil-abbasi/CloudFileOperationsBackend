package files

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	errs "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/errors"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FilesService struct {
	repository interfaces.IFilesRepository
	config     *config.Config
}

func NewService(config *config.Config) *FilesService {

	repository := NewSQLiteRepository(config)

	return &FilesService{
		repository: repository,
		config:     config,
	}
}

func (s *FilesService) FindOne(id string) (entities.File, bool, error) {
	return s.repository.FindOne(id)
}

func (s *FilesService) Create(dto dtos.FileCreateDto) (entities.File, error) {
	filePath := s.config.RootPath + dto.Location + "/" + dto.Filename
	exists, err := shared.DoesFileExist(filePath)

	if err != nil {
		return entities.File{}, err
	}

	if exists {
		return entities.File{}, &errs.FileAlreadyExistsError{}
	}

	return s.repository.Create(dto)
}

func (s *FilesService) Update(dto dtos.FileUpdateDto) (entities.File, bool, error) {
	return s.repository.Update(dto)
}

func (s *FilesService) Remove(id string) (entities.File, bool, error) {
	return s.repository.Remove(id)
}
