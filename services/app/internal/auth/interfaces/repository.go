package interfaces

import "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"

type IUsersRepository interface {
	FindByName(name string) (entities.User, bool, error)
	Save(entity entities.User) error
	Remove(id string) (bool, error)
}
