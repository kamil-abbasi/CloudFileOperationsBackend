package directories

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/pkg"
)

type DirectoriesService struct {
	directoriesRepository IDirectoriesRepository
	config                *config.Config
	filesRepository       IFilesRepository
	storage               storage.IStorage
}

func NewService(config *config.Config, directoriesRepository IDirectoriesRepository, filesRepository IFilesRepository, storage storage.IStorage) *DirectoriesService {
	return &DirectoriesService{
		directoriesRepository: directoriesRepository,
		config:                config,
		filesRepository:       filesRepository,
		storage:               storage,
	}
}

func (s *DirectoriesService) ListItems(location string) ([]dtos.DirectoryItemResponseDto, error) {
	items, err := s.directoriesRepository.ListItems(location)

	if err != nil {
		return []dtos.DirectoryItemResponseDto{}, err
	}

	result := pkg.SliceMap(items, ItemEntityToDto)

	return result, nil
}

func (s *DirectoriesService) FindOne(id string) (dtos.DirectoryResponseDto, error) {
	dir, found, err := s.directoriesRepository.FindOne(id)

	if err != nil {
		return dtos.DirectoryResponseDto{}, err
	}

	if !found {
		return dtos.DirectoryResponseDto{}, ErrNotFound
	}

	return EntityToDto(dir), nil
}

func (s *DirectoriesService) Create(userId string, dto dtos.CreateDirectoryDto) (dtos.DirectoryResponseDto, error) {
	var parentId = ""
	var location = "/"

	if dto.ParentId != "" {
		_, found, err := s.directoriesRepository.FindByNameAndParentId(dto.Name, dto.ParentId)

		if err != nil {
			return dtos.DirectoryResponseDto{}, err
		}

		if found {
			return dtos.DirectoryResponseDto{}, ErrAlreadyExists
		}

		parentDir, found, err := s.directoriesRepository.FindOne(dto.ParentId)

		if err != nil {
			return dtos.DirectoryResponseDto{}, err
		}

		if !found {
			return dtos.DirectoryResponseDto{}, ErrNotFound
		}

		parentId = parentDir.Id
		location = filepath.Join(parentDir.Location, parentDir.Name)
	}

	directory := entities.Directory{
		Id:       uuid.NewString(),
		UserId:   userId,
		ParentId: parentId,
		Name:     dto.Name,
		Location: filepath.Clean(location),
	}

	err := s.directoriesRepository.Save(directory)

	if err != nil {
		return dtos.DirectoryResponseDto{}, err
	}

	return EntityToDto(directory), nil
}

// not implemented
func (s *DirectoriesService) Update(id string, dto dtos.UpdateDirectoryDto) (dtos.DirectoryResponseDto, error) {
	return dtos.DirectoryResponseDto{}, fmt.Errorf("operation not implemented")
}

// not implemented
func (s *DirectoriesService) Rename(id string, newName string) (dtos.DirectoryResponseDto, error) {
	return dtos.DirectoryResponseDto{}, fmt.Errorf("operation not implemented")
}

// not implemented
func (s *DirectoriesService) Move(id string, parentId string) (dtos.DirectoryResponseDto, error) {
	return dtos.DirectoryResponseDto{}, fmt.Errorf("operation not implemented")
}

func (s *DirectoriesService) Remove(id string) error {
	dir, found, err := s.directoriesRepository.FindOne(id)

	if err != nil {
		return err
	}

	if !found {
		return ErrNotFound
	}

	files, err := s.filesRepository.FindByLocation(filepath.Join(dir.Location, dir.Name))

	if err != nil {
		return err
	}

	for _, file := range files {
		s.storage.RemoveFile(file.Id)
	}

	_, err = s.directoriesRepository.Remove(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *DirectoriesService) Download(id string, writer io.Writer) error {
	directory, found, err := s.directoriesRepository.FindOne(id)

	if err != nil {
		return err
	}

	if !found {
		return ErrNotFound
	}

	files, err := s.filesRepository.FindByLocation(filepath.Join(directory.Location, directory.Name))

	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// base path for every file in the directory
	basePath := filepath.Join(directory.Location, directory.Name) + string(filepath.Separator)

	for _, file := range files {

		fullFilePath := filepath.Join(file.Location, file.Name)
		relativePath := strings.TrimPrefix(fullFilePath, basePath)
		relativePath = filepath.ToSlash(relativePath)

		f, err := zipWriter.Create(relativePath)

		// TODO: handle error properly
		if err != nil {
			return err
		}

		readCloser, found, err := s.storage.DownloadFile(file.Id)

		// TODO: handle error properly
		if err != nil || !found {
			return err
		}

		_, err = io.Copy(f, readCloser)

		readCloser.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
