package files

import (
	"archive/zip"
	"fmt"
	"io"
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

func (s *FilesService) CreateDirectory(userId string, location string) error {
	dirPath := filepath.Clean(
		filepath.Join(s.config.RootPath, userId, location),
	)

	return os.MkdirAll(dirPath, 0755)
}

func (s *FilesService) FindByLocation(userId string, location string) ([]entities.File, error) {
	return s.repository.FindByLocation(userId, location)
}

func (s *FilesService) RemoveDirectory(userId string, location string) (int64, error) {
	dirPath := filepath.Clean(
		filepath.Join(s.config.RootPath, userId, location),
	)

	err := os.RemoveAll(dirPath)

	if err != nil {
		return 0, fmt.Errorf("error while removing directory: %v", err)
	}

	return s.repository.RemoveByLocation(userId, location)
}

func (s *FilesService) DownloadDirectory(userId string, location string, writer io.Writer) error {
	files, err := s.repository.FindByLocation(userId, location)

	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	for _, file := range files {
		fullPath := filepath.Clean(
			filepath.Join(s.config.RootPath, file.UserId, file.Location, file.Filename),
		)

		// relative path inside ZIP based on location prefix
		relativePath := filepath.Join(file.Location, file.Filename)
		if len(location) > 0 {
			rel, err := filepath.Rel(location, relativePath)
			if err == nil {
				relativePath = rel
			}
		}

		srcFile, err := os.Open(fullPath)

		if err != nil {
			return fmt.Errorf("error while opening file %v: %v", fullPath, err)
		}

		zipEntry, err := zipWriter.Create(relativePath)

		if err != nil {
			srcFile.Close()
			return fmt.Errorf("error while creating zip entry: %v", err)
		}

		_, err = io.Copy(zipEntry, srcFile)
		srcFile.Close()

		if err != nil {
			return fmt.Errorf("error while writing to zip: %v", err)
		}
	}

	return nil
}
