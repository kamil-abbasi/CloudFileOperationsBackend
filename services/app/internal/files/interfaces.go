package files

import (
	usersEntities "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"
	directoriesEntities "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

type IFilesRepository interface {
	Find(dto dtos.FindFilesDto) ([]entities.File, error)
	FindByLocation(location string) ([]entities.File, error)
	FindOne(id string) (entities.File, bool, error)
	FindByNameAndDirectoryId(name string, directoryId string) (entities.File, bool, error)
	Save(file entities.File) error
	Remove(id string) (bool, error)
}

type IDirectoriesRepository interface {
	FindOne(id string) (directoriesEntities.Directory, bool, error)
}

type IUsersRepository interface {
	FindOne(id string) (usersEntities.User, bool, error)
	Save(user usersEntities.User) error
}
