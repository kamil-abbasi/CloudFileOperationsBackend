package files

import (
	"io"

	"github.com/google/uuid"
	directoriesInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/helpers"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FilesService struct {
	repository            interfaces.IFilesRepository
	storageAdapter        shared.IStorageAdapter
	directoriesRepository directoriesInterfaces.IDirectoriesRepository
}

func NewService(repository interfaces.IFilesRepository, directoriesRepository directoriesInterfaces.IDirectoriesRepository, storageAdapter shared.IStorageAdapter) FilesService {
	return FilesService{
		repository:            repository,
		directoriesRepository: directoriesRepository,
		storageAdapter:        storageAdapter,
	}
}

func (s *FilesService) FindOne(id string) (entities.File, bool, error) {
	file, found, err := s.repository.FindOne(id)

	if err != nil {
		return entities.File{}, false, err
	}

	if !found {
		return entities.File{}, false, nil
	}

	return file, true, nil
}

func (s *FilesService) Download(id string) (io.ReadCloser, bool, error) {
	file, found, err := s.repository.FindOne(id)

	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	key := helpers.GetFileKey(file)

	src, err := s.storageAdapter.DownloadFile(key)

	return src, true, nil
}

/*
exceptions to handle:
shared.DirectoryNotFound
shared.FileAlreadyExists
*/
func (s *FilesService) Create(dto dtos.CreateFileDto, reader io.Reader) (entities.File, error) {
	// if directory id is empty file will be in the root directory
	if dto.DirectoryId != "" {
		_, found, err := s.directoriesRepository.FindOne(dto.DirectoryId)

		if err != nil {
			return entities.File{}, err
		}

		if !found {
			return entities.File{}, &shared.DirectoryNotFoundError{}
		}
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
	}

	key := helpers.GetFileKey(file)

	bytesWritten, err := s.storageAdapter.UploadFile(key, reader)

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
		_, found, err = s.directoriesRepository.FindOne(file.DirectoryId)

		if err != nil || !found {
			return entities.File{}, false, err
		}

		file.DirectoryId = dto.Fields.DirectoryId
	}

	file.Name = dto.Fields.Name

	s.repository.Save(file)

	return file, true, nil
}

func (s *FilesService) Remove(id string) (bool, error) {
	file, found, err := s.repository.FindOne(id)

	if err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	key := helpers.GetFileKey(file)

	wasRemoved, err := s.storageAdapter.RemoveFile(key)

	if err != nil || !wasRemoved {
		return false, err
	}

	found, err = s.repository.Remove(id)

	if err != nil || !found {
		return false, err
	}

	return true, nil
}
