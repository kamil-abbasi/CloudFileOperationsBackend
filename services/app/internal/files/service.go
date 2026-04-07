package files

import (
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	directoriesInterfaces "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
)

type FilesService struct {
	repository            interfaces.IFilesRepository
	storage               storage.IStorage
	directoriesRepository directoriesInterfaces.IDirectoriesRepository
}

func NewService(repository interfaces.IFilesRepository, directoriesRepository directoriesInterfaces.IDirectoriesRepository, storage storage.IStorage) FilesService {
	return FilesService{
		repository:            repository,
		directoriesRepository: directoriesRepository,
		storage:               storage,
	}
}

func (s *FilesService) FindByLocation(location string) ([]entities.File, error) {
	return s.repository.FindByLocation(location)
}

func (s *FilesService) FindOne(id string) (entities.File, error) {
	file, found, err := s.repository.FindOne(id)

	if err != nil {
		return entities.File{}, err
	}

	if !found {
		return entities.File{}, ErrNotFound
	}

	return file, nil
}

func (s *FilesService) Download(id string) (io.ReadCloser, error) {
	reader, found, err := s.storage.DownloadFile(id)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrNotFound
	}

	return reader, nil
}

/*
exceptions to handle:
shared.DirectoryNotFound
shared.FileAlreadyExists
*/
func (s *FilesService) Create(userId string, dto dtos.CreateFileDto, reader io.Reader) (entities.File, error) {
	// if directory id is empty file will be in the root directory
	location := "/"

	if dto.DirectoryId != "" {
		parentDir, found, err := s.directoriesRepository.FindOne(dto.DirectoryId)

		if err != nil {
			return entities.File{}, err
		}

		if !found {
			return entities.File{}, directories.ErrNotFound
		}

		location = filepath.Clean(filepath.Join(parentDir.Location, parentDir.Name))
	}

	_, found, err := s.repository.FindByNameAndDirectoryId(dto.Name, dto.DirectoryId)

	if found {
		return entities.File{}, ErrAlreadyExists
	}

	file := entities.File{
		Id:          uuid.NewString(),
		UserId:      userId,
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

func (s *FilesService) Update(id string, dto dtos.UpdateFileDto) (entities.File, error) {
	file, found, err := s.repository.FindOne(id)

	if dto.Name == "" {
		dto.Name = file.Name
	}

	if err != nil {
		return entities.File{}, err
	}

	if !found {
		return entities.File{}, ErrNotFound
	}

	if dto.DirectoryId != "" {
		dir, found, err := s.directoriesRepository.FindOne(file.DirectoryId)

		if err != nil || !found {
			return entities.File{}, err
		}

		file.DirectoryId = dto.DirectoryId
		file.Location = filepath.Join(dir.Location, dir.Name)
	}

	file.Name = dto.Name

	s.repository.Save(file)

	return file, nil
}

func (s *FilesService) Remove(id string) error {
	wasRemoved, err := s.repository.Remove(id)

	if err != nil {
		return err
	}

	if !wasRemoved {
		return ErrNotFound
	}

	removed, err := s.storage.RemoveFile(id)

	if !removed || err != nil {
		return err
	}

	return nil
}
