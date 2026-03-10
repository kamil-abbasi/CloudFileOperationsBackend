package directories

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	fileDtos "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	fileInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type DirectoriesService struct {
	repository      interfaces.IDirectoriesRepository
	config          *config.Config
	filesRepository fileInterfaces.IFilesRepository
}

func NewService(config *config.Config, repository interfaces.IDirectoriesRepository, filesRepository fileInterfaces.IFilesRepository) DirectoriesService {
	return DirectoriesService{
		repository:      repository,
		config:          config,
		filesRepository: filesRepository,
	}
}

func (s *DirectoriesService) Create(dto dtos.CreateDirectoryDto) (entities.Directory, error) {
	path := filepath.Clean(filepath.Join(s.config.RootPath, dto.UserId, dto.Location))

	err := os.MkdirAll(path, 0755)

	if err != nil {
		return entities.Directory{}, err
	}

	return entities.Directory(dto), nil
}

func (s *DirectoriesService) Download(id string, writer io.Writer) error {
	directory, found, err := s.repository.FindOne(id)

	if err != nil {
		return err
	}

	if !found {
		return nil
	}

	findDto := fileDtos.FindFilesDto{}
	findDto.Where.DirectoryId = id

	files, err := s.filesRepository.Find(findDto)

	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	for _, file := range files {
		fullPath := filepath.Clean(filepath.Join(s.config.RootPath, file.UserId, file.Location, file.Name))

		// relative path inside ZIP based on location prefix
		relativePath := filepath.Join(file.Location, file.Name)

		if len(directory.Location) > 0 {
			rel, err := filepath.Rel(directory.Location, relativePath)

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

func (s *DirectoriesService) Remove(id string) (bool, error) {
	directory, found, err := s.repository.FindOne(id)

	if !found {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	path := filepath.Clean(filepath.Join(s.config.RootPath, directory.UserId, directory.Location))

	err = os.RemoveAll(path)

	if err != nil {
		return false, err
	}

	s.repository.Remove(id)

	return true, nil
}
