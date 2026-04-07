package auth

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/interfaces"
)

type PostgresUsersRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) interfaces.IUsersRepository {
	return &PostgresUsersRepository{
		db: db,
	}
}

func (r *PostgresUsersRepository) FindByName(name string) (entities.User, bool, error) {
	if name == "" {
		return entities.User{}, false, nil
	}

	var user entities.User

	err := r.db.QueryRow(`
		SELECT id, name, password_hash
		FROM users
		WHERE name = $1
	`, name).Scan(&user.Id, &user.Name, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return entities.User{}, false, nil
	}

	if err != nil {
		return entities.User{}, false, fmt.Errorf("failed to find user by name, details: %v", err)
	}

	return user, true, nil
}

func (r *PostgresUsersRepository) Save(entity entities.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users(id, name, password_hash)
		VALUES($1, $2, $3)
		ON CONFLICT(id)
		DO UPDATE SET
			name = EXCLUDED.name
	`, entity.Id, entity.Name, entity.PasswordHash)

	if err != nil {
		return fmt.Errorf("failed to save user to postgres, details: %v", err)
	}

	return nil
}

func (r *PostgresUsersRepository) Remove(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)

	if err != nil {
		return false, fmt.Errorf("failed to remove user from postgres, details: %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, fmt.Errorf("failed to fetch the number of deleted rows: %v", err)
	}

	if rowsAffected <= 0 {
		return false, nil
	}

	return true, nil
}
