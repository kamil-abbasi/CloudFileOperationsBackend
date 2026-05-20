package files

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/storage"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/pkg"
)

type FilesService struct {
	filesRepository       IFilesRepository
	storage               storage.IStorage
	directoriesRepository IDirectoriesRepository
	usersRepository       IUsersRepository
}

func NewService(filesRepository IFilesRepository, directoriesRepository IDirectoriesRepository, storage storage.IStorage, usersRepository IUsersRepository) *FilesService {
	return &FilesService{
		filesRepository:       filesRepository,
		directoriesRepository: directoriesRepository,
		storage:               storage,
		usersRepository:       usersRepository,
	}
}

func (s *FilesService) FindByLocation(location string) ([]dtos.FileResponseDto, error) {
	entities, err := s.filesRepository.FindByLocation(location)

	if err != nil {
		return []dtos.FileResponseDto{}, err
	}

	result := pkg.SliceMap(entities, EntityToDto)

	return result, nil
}

func (s *FilesService) FindOne(id string) (dtos.FileResponseDto, error) {
	file, found, err := s.filesRepository.FindOne(id)

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	if !found {
		return dtos.FileResponseDto{}, ErrNotFound
	}

	return EntityToDto(file), nil
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
func (s *FilesService) Create(userId string, dto dtos.CreateFileDto, reader io.Reader) (dtos.FileResponseDto, error) {
	// if directory id is empty file will be in the root directory
	location := "/"

	user, found, err := s.usersRepository.FindOne(userId)

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	if !found {
		return dtos.FileResponseDto{}, ErrUserNotFound
	}

	if dto.DirectoryId != "" {
		parentDir, found, err := s.directoriesRepository.FindOne(dto.DirectoryId)

		if err != nil {
			return dtos.FileResponseDto{}, err
		}

		if !found {
			return dtos.FileResponseDto{}, ErrDirectoryNotFound
		}

		location = filepath.Clean(filepath.Join(parentDir.Location, parentDir.Name))
	}

	_, found, err = s.filesRepository.FindByNameAndDirectoryId(dto.Name, dto.DirectoryId)

	if found {
		return dtos.FileResponseDto{}, ErrAlreadyExists
	}

	fileId := uuid.NewString()

	bytesWritten, checksum, err := s.storage.UploadFile(fileId, reader)

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	file := entities.File{
		Id:          fileId,
		UserId:      userId,
		DirectoryId: dto.DirectoryId,
		Name:        dto.Name,
		Size:        0,
		Location:    location,
	}

	file.Size = uint64(bytesWritten)
	file.Checksum = checksum

	user.StorageUsed += file.Size

	if user.StorageUsed > user.MaxStorage {
		fmt.Print(user)
		s.storage.RemoveFile(file.Id)

		return dtos.FileResponseDto{}, ErrNotEnoughStorage
	}

	if file.Checksum != dto.Checksum {
		return dtos.FileResponseDto{}, ErrCorruptedUpload
	}

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	err = s.filesRepository.Save(file)

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	err = s.usersRepository.Save(user)

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	return EntityToDto(file), nil
}

func (s *FilesService) Update(id string, dto dtos.UpdateFileDto) (dtos.FileResponseDto, error) {
	file, found, err := s.filesRepository.FindOne(id)

	if dto.Name == "" {
		dto.Name = file.Name
	}

	if err != nil {
		return dtos.FileResponseDto{}, err
	}

	if !found {
		return dtos.FileResponseDto{}, ErrNotFound
	}

	if dto.DirectoryId != "" {
		dir, found, err := s.directoriesRepository.FindOne(file.DirectoryId)

		if err != nil || !found {
			return dtos.FileResponseDto{}, err
		}

		file.DirectoryId = dto.DirectoryId
		file.Location = filepath.Join(dir.Location, dir.Name)
	}

	file.Name = dto.Name

	s.filesRepository.Save(file)

	return EntityToDto(file), nil
}

func (s *FilesService) Remove(id string) error {
	wasRemoved, err := s.filesRepository.Remove(id)

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
