package directories

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	fileInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type DirectoriesService struct {
	repository      interfaces.IDirectoriesRepository
	config          *config.Config
	filesRepository fileInterfaces.IFilesRepository
	storage         shared.IStorage
}

func NewService(config *config.Config, repository interfaces.IDirectoriesRepository, filesRepository fileInterfaces.IFilesRepository, storage shared.IStorage) DirectoriesService {
	return DirectoriesService{
		repository:      repository,
		config:          config,
		filesRepository: filesRepository,
		storage:         storage,
	}
}

func (s *DirectoriesService) FindOne(id string) (entities.Directory, bool, error) {
	return s.repository.FindOne(id)
}

func (s *DirectoriesService) Create(dto dtos.CreateDirectoryDto) (entities.Directory, error) {
	var parentId = ""
	var location = "/"

	if dto.ParentId != "" {
		_, found, err := s.repository.FindByNameAndParentId(dto.Name, dto.ParentId)

		if err != nil {
			return entities.Directory{}, err
		}

		if found {
			return entities.Directory{}, &shared.DirectoryAlreadyExistsError{}
		}

		parentDir, found, err := s.repository.FindOne(dto.ParentId)

		if err != nil {
			return entities.Directory{}, err
		}

		if !found {
			return entities.Directory{}, &shared.DirectoryNotFoundError{}
		}

		parentId = parentDir.Id
		location = filepath.Join(parentDir.Location, parentDir.Name)
	}

	log.Printf("location: %v", location)

	directory := entities.Directory{
		Id:       uuid.NewString(),
		UserId:   dto.UserId,
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
func (s *DirectoriesService) Update(dto dtos.UpdateDirectoryDto) (entities.Directory, bool, error) {
	return entities.Directory{}, false, fmt.Errorf("operation not implemented")
}

func (s *DirectoriesService) Remove(id string) (bool, error) {
	dir, found, err := s.repository.FindOne(id)

	if err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	files, err := s.filesRepository.FindByLocation(filepath.Join(dir.Location, dir.Name))

	if err != nil {
		return false, err
	}

	for _, file := range files {
		s.storage.RemoveFile(file.Id)
	}

	_, err = s.repository.Remove(id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *DirectoriesService) Download(id string, writer io.Writer) (bool, error) {
	directory, found, err := s.repository.FindOne(id)

	if err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	files, err := s.filesRepository.FindByLocation(filepath.Join(directory.Location, directory.Name))

	if err != nil {
		return false, err
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
			return false, err
		}

		readCloser, found, err := s.storage.DownloadFile(file.Id)

		// TODO: handle error properly
		if err != nil || !found {
			return false, err
		}

		_, err = io.Copy(f, readCloser)

		readCloser.Close()

		if err != nil {
			return false, err
		}
	}

	return true, nil
}
