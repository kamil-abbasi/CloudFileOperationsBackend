package files

import (
	"io"
	"path/filepath"

	"github.com/google/uuid"
	directoriesInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FilesService struct {
	repository            interfaces.IFilesRepository
	storage               shared.IStorage
	directoriesRepository directoriesInterfaces.IDirectoriesRepository
}

func NewService(repository interfaces.IFilesRepository, directoriesRepository directoriesInterfaces.IDirectoriesRepository, storage shared.IStorage) FilesService {
	return FilesService{
		repository:            repository,
		directoriesRepository: directoriesRepository,
		storage:               storage,
	}
}

func (s *FilesService) FindByLocation(location string) ([]entities.File, error) {
	return s.repository.FindByLocation(location)
}

func (s *FilesService) FindOne(id string) (entities.File, bool, error) {
	return s.repository.FindOne(id)
}

func (s *FilesService) Download(id string) (io.ReadCloser, bool, error) {
	return s.storage.DownloadFile(id)
}

/*
exceptions to handle:
shared.DirectoryNotFound
shared.FileAlreadyExists
*/
func (s *FilesService) Create(dto dtos.CreateFileDto, reader io.Reader) (entities.File, error) {
	// if directory id is empty file will be in the root directory
	location := "/"

	if dto.DirectoryId != "" {
		parentDir, found, err := s.directoriesRepository.FindOne(dto.DirectoryId)

		if err != nil {
			return entities.File{}, err
		}

		if !found {
			return entities.File{}, &shared.DirectoryNotFoundError{}
		}

		location = filepath.Clean(filepath.Join(parentDir.Location, parentDir.Name))
	}

	_, found, err := s.repository.FindByNameAndDirectoryId(dto.Name, dto.DirectoryId)

	if found {
		return entities.File{}, &shared.FileAlreadyExistsError{}
	}

	file := entities.File{
		Id:          uuid.NewString(),
		UserId:      dto.UserId,
		DirectoryId: dto.DirectoryId,
		Name:        dto.Name,
		Size:        0,
		Location:    location,
	}

	bytesWritten, err := s.storage.UploadFile(file.Id, reader)

	file.Size = uint64(bytesWritten)

	if err != nil {
		return entities.File{}, err
	}

	err = s.repository.Save(file)

	if err != nil {
		return entities.File{}, err
	}

	return file, nil
}

func (s *FilesService) Update(dto dtos.UpdateFileDto) (entities.File, bool, error) {
	file, found, err := s.repository.FindOne(dto.Where.Id)

	if dto.Fields.Name == "" {
		dto.Fields.Name = file.Name
	}

	if err != nil {
		return entities.File{}, false, err
	}

	if !found {
		return entities.File{}, false, nil
	}

	if dto.Fields.DirectoryId != "" {
		dir, found, err := s.directoriesRepository.FindOne(file.DirectoryId)

		if err != nil || !found {
			return entities.File{}, false, err
		}

		file.DirectoryId = dto.Fields.DirectoryId
		file.Location = filepath.Join(dir.Location, dir.Name)
	}

	file.Name = dto.Fields.Name

	s.repository.Save(file)

	return file, true, nil
}

func (s *FilesService) Remove(id string) (bool, error) {
	wasRemoved, err := s.repository.Remove(id)

	if err != nil {
		return false, err
	}

	if !wasRemoved {
		return false, nil
	}

	removed, err := s.storage.RemoveFile(id)

	if !removed || err != nil {
		return false, err
	}

	return true, nil
}
