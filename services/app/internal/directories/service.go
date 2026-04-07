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
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	fileInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
)

type DirectoriesService struct {
	repository      interfaces.IDirectoriesRepository
	config          *config.Config
	filesRepository fileInterfaces.IFilesRepository
	storage         storage.IStorage
}

func NewService(config *config.Config, repository interfaces.IDirectoriesRepository, filesRepository fileInterfaces.IFilesRepository, storage storage.IStorage) DirectoriesService {
	return DirectoriesService{
		repository:      repository,
		config:          config,
		filesRepository: filesRepository,
		storage:         storage,
	}
}

func (s *DirectoriesService) ListItems(location string) ([]entities.DirectoryItem, error) {
	return s.repository.ListItems(location)
}

func (s *DirectoriesService) FindOne(id string) (entities.Directory, error) {
	dir, found, err := s.repository.FindOne(id)

	if err != nil {
		return entities.Directory{}, err
	}

	if !found {
		return entities.Directory{}, ErrNotFound
	}

	return dir, nil
}

func (s *DirectoriesService) Create(userId string, dto dtos.CreateDirectoryDto) (entities.Directory, error) {
	var parentId = ""
	var location = "/"

	if dto.ParentId != "" {
		_, found, err := s.repository.FindByNameAndParentId(dto.Name, dto.ParentId)

		if err != nil {
			return entities.Directory{}, err
		}

		if found {
			return entities.Directory{}, ErrAlreadyExists
		}

		parentDir, found, err := s.repository.FindOne(dto.ParentId)

		if err != nil {
			return entities.Directory{}, err
		}

		if !found {
			return entities.Directory{}, ErrNotFound
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

	err := s.repository.Save(directory)

	if err != nil {
		return entities.Directory{}, err
	}

	return directory, nil
}

// not implemented
func (s *DirectoriesService) Update(id string, dto dtos.UpdateDirectoryDto) (entities.Directory, error) {
	return entities.Directory{}, fmt.Errorf("operation not implemented")
}

// not implemented
func (s *DirectoriesService) Rename(id string, newName string) (entities.Directory, error) {
	return entities.Directory{}, fmt.Errorf("operation not implemented")
}

// not implemented
func (s *DirectoriesService) Move(id string, parentId string) (entities.Directory, error) {
	return entities.Directory{}, fmt.Errorf("operation not implemented")
}

func (s *DirectoriesService) Remove(id string) error {
	dir, found, err := s.repository.FindOne(id)

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

	_, err = s.repository.Remove(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *DirectoriesService) Download(id string, writer io.Writer) error {
	directory, found, err := s.repository.FindOne(id)

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
