package interfaces

import "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"

type IUsersRepository interface {
	FindOne(id string) (entities.User, bool, error)
	FindByName(name string) (entities.User, bool, error)
	Save(entity entities.User) error
	Remove(id string) (bool, error)
}
