package auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/interfaces"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
)

type PostgresUsersRepository struct {
	db *database.Postgres
}

func NewPostgresRepository(db *database.Postgres) interfaces.IUsersRepository {
	return &PostgresUsersRepository{
		db: db,
	}
}

func (r *PostgresUsersRepository) FindByName(name string) (entities.User, bool, error) {
	if name == "" {
		return entities.User{}, false, nil
	}

	user, err := r.db.Queries.FindUserByName(context.Background(), name)

	if err == sql.ErrNoRows {
		return entities.User{}, false, nil
	}

	if err != nil {
		return entities.User{}, false, fmt.Errorf("failed to find user by name, details: %v", err)
	}

	return entities.User{
		Id:           user.ID.String(),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		MaxStorage:   uint64(user.MaxStorage),
		StorageUsed:  uint64(user.StorageUsed),
	}, true, nil
}

func (r *PostgresUsersRepository) FindOne(id string) (entities.User, bool, error) {
	if id == "" {
		return entities.User{}, false, nil
	}

	userId, err := uuid.Parse(id)

	if err != nil {
		return entities.User{}, false, fmt.Errorf("failed to parse user id, details: %v", err)
	}

	user, err := r.db.Queries.FindUserById(context.Background(), userId)

	if err == sql.ErrNoRows {
		return entities.User{}, false, nil
	}

	if err != nil {
		return entities.User{}, false, fmt.Errorf("failed to find user by id, details: %v", err)
	}

	return entities.User{
		Id:           user.ID.String(),
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		MaxStorage:   uint64(user.MaxStorage),
		StorageUsed:  uint64(user.StorageUsed),
	}, true, nil
}

func (r *PostgresUsersRepository) Save(entity entities.User) error {
	userID, err := uuid.Parse(entity.Id)

	if err != nil {
		return fmt.Errorf("failed to parse user id, details: %v", err)
	}

	err = r.db.Queries.SaveUser(context.Background(), database.SaveUserParams{
		ID:           userID,
		Name:         entity.Name,
		PasswordHash: entity.PasswordHash,
		MaxStorage:   int64(entity.MaxStorage),
		StorageUsed:  int64(entity.StorageUsed),
	})

	if err != nil {
		return fmt.Errorf("failed to save user to postgres, details: %v", err)
	}

	return nil
}

func (r *PostgresUsersRepository) Remove(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	userID, err := uuid.Parse(id)

	if err != nil {
		return false, fmt.Errorf("failed to parse user id, details: %v", err)
	}

	rowsAffected, err := r.db.Queries.RemoveUserRows(context.Background(), userID)

	if err != nil {
		return false, fmt.Errorf("failed to remove user from postgres, details: %v", err)
	}

	if rowsAffected <= 0 {
		return false, nil
	}

	return true, nil
}
