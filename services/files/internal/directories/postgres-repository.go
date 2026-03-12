package directories

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/interfaces"
)

type DirectoriesPostgresRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewPostgresRepository(config *config.Config, db *sql.DB) interfaces.IDirectoriesRepository {
	return &DirectoriesPostgresRepository{
		db:     db,
		config: config,
	}
}

func (repository *DirectoriesPostgresRepository) Find() ([]entities.Directory, error) {
	var directories []entities.Directory

	rows, err := repository.db.Query(`
		SELECT id, user_id, parent_id, name
		FROM directories`)

	if err != nil {
		return nil, fmt.Errorf("Failed to find directories: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var directory entities.Directory
		var parentId sql.NullString

		err := rows.Scan(&directory.Id, &directory.UserId, &parentId, &directory.Name)

		if err != nil {
			return nil, err
		}

		if parentId.Valid {
			directory.ParentId = parentId.String
		} else {
			directory.ParentId = ""
		}

		directories = append(directories, directory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return directories, nil
}

func (repository *DirectoriesPostgresRepository) FindOne(id string) (entities.Directory, bool, error) {
	if id == "" {
		return entities.Directory{}, false, nil
	}

	var directory entities.Directory
	var parentId sql.NullString

	err := repository.db.QueryRow(`
		SELECT id, user_id, parent_id, name
		FROM directories
		WHERE id = $1`,
		id).Scan(&directory.Id, &directory.UserId, &parentId, &directory.Name)

	if err == sql.ErrNoRows {
		return entities.Directory{}, false, nil
	}

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("Failed to find one file. Details: %v", err)
	}

	if parentId.Valid {
		directory.ParentId = parentId.String
	} else {
		directory.ParentId = ""
	}

	return directory, true, nil
}

func (repository *DirectoriesPostgresRepository) Save(directory entities.Directory) error {
	_, err := repository.db.Exec(
		`INSERT INTO
		directories(id, user_id, parent_id, name)
		VALUES($1, $2, $3, $4)
		ON CONFLICT(id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			parent_id = EXCLUDED.parent_id,
			name = EXCLUDED.name
		RETURNING id, user_id, parent_id, name`,
		directory.Id, directory.UserId, &sql.NullString{
			String: directory.ParentId,
		}, directory.Name)

	if err != nil {
		return fmt.Errorf("failed to save directory to postgres, details: %v", err)
	}

	return nil
}

func (repository *DirectoriesPostgresRepository) Remove(id string) (bool, error) {
	result, err := repository.db.Exec("DELETE FROM directories WHERE id = $1", id)

	if err != nil {
		return false, fmt.Errorf("failed to remove directory from postgres, details: %v", err)
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
