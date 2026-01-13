package files

import (
	"fmt"
	"os"
	"path/filepath"

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
	fullPath := filepath.Clean(
		filepath.Join(
			s.config.RootPath,
			dto.UserId,
			dto.Location,
			dto.Filename,
		),
	)

	fmt.Println(fullPath)

	exists, err := shared.DoesFileExist(fullPath)

	if err != nil {
		return entities.File{}, err
	}

	if exists {
		return entities.File{}, &errs.FileAlreadyExistsError{}
	}

	return s.repository.Create(dto)
}

func (s *FilesService) CreateFolder(dto dtos.FolderCreateDto) (entities.Folder, error) {
	fullPath := filepath.Clean(
		filepath.Join(
			s.config.RootPath,
			dto.UserId,
			dto.Location,
			dto.Name,
		),
	)

	exists, err := shared.DoesFileExist(fullPath)

	if err != nil {
		return entities.Folder{}, err
	}

	if exists {
		return entities.Folder{}, &errs.FolderAlreadyExistsError{}
	}

	err = os.MkdirAll(fullPath, 0755)

	if err != nil {
		return entities.Folder{}, fmt.Errorf("error while creating directory: %v", err)
	}

	return s.repository.CreateFolder(dto)
}

func (s *FilesService) Update(dto dtos.FileUpdateDto) (entities.File, bool, error) {
	file, found, err := s.repository.FindOne(dto.Where.Id)
	fullPath := filepath.Join(filepath.Join(
		s.config.RootPath,
		file.UserId,
		file.Location,
		file.Filename,
	))

	if err != nil {
		return entities.File{}, false, err
	}

	if !found {
		return entities.File{}, false, nil
	}

	if dto.Fields.Filename == "" {
		dto.Fields.Filename = file.Filename
	}

	if dto.Fields.Location == "" {
		dto.Fields.Location = file.Location
	}

	fmt.Println(dto)

	newFullPath := filepath.Join(filepath.Join(
		s.config.RootPath,
		file.UserId,
		dto.Fields.Location,
		dto.Fields.Filename,
	))

	err = os.MkdirAll(filepath.Dir(newFullPath), 0755)

	if err != nil {
		return entities.File{}, false, fmt.Errorf("error while creating dirs: %v", err)
	}

	err = os.Rename(fullPath, newFullPath)

	if err != nil {
		return entities.File{}, false, fmt.Errorf("error while moving file: %v", err)
	}

	s.repository.Update(dto)

	return entities.File{
		Id:       file.Id,
		UserId:   file.UserId,
		Size:     file.Size,
		Location: dto.Fields.Location,
		Filename: dto.Fields.Filename,
	}, true, nil
}

func (s *FilesService) Remove(id string) (entities.File, bool, error) {
	file, found, err := s.repository.FindOne(id)

	if err != nil {
		return entities.File{}, false, err
	}

	if !found {
		return entities.File{}, false, nil
	}

	fullPath := filepath.Clean(
		filepath.Join(
			s.config.RootPath,
			file.UserId,
			file.Location,
			file.Filename,
		),
	)

	err = os.Remove(fullPath)

	if err != nil {
		return entities.File{}, false, err
	}

	s.repository.Remove(id)

	return file, true, nil
}
