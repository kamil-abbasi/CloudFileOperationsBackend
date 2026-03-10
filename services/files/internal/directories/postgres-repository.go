package directories

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
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

	rows, err := repository.db.Query("SELECT * FROM directories")

	if err != nil {
		return nil, fmt.Errorf("Failed to find directories: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var directory entities.Directory

		err := rows.Scan(directory.Id, directory.Name, directory.Location, directory.UserId, directory.ParentId)

		if err != nil {
			return nil, err
		}

		directories = append(directories, directory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return directories, nil
}

func (repository *DirectoriesPostgresRepository) FindOne(id string) (entities.Directory, bool, error) {
	var directory entities.Directory

	err := repository.db.QueryRow("SELECT * FROM directories WHERE id = $1", id).Scan(&directory.Id, &directory.Name, &directory.Location, &directory.UserId, &directory.ParentId)

	if err == sql.ErrNoRows {
		return entities.Directory{}, false, nil
	}

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("Failed to find one directory. Details: %v", err)
	}

	return directory, true, nil
}

func (repository *DirectoriesPostgresRepository) Create(dto dtos.CreateDirectoryDto) (entities.Directory, error) {
	_, err := repository.db.Exec(
		`INSERT INTO
		directories(id, name, location, user_id, parent_id)
		VALUES($1, $2, $3, $4, $5, $6)`,
		dto.Id, dto.Name, dto.Location, dto.UserId, dto.UserId, dto.ParentId)

	if err != nil {
		return entities.Directory{}, fmt.Errorf("failed to create directory metadata, details: %v", err)
	}

	return entities.Directory(dto), nil
}

func (repository *DirectoriesPostgresRepository) Remove(id string) (bool, error) {
	result, err := repository.db.Exec("DELETE FROM directories WHERE id = $1", id)

	if err != nil {
		return false, fmt.Errorf("failed to remove directory metadata, details: %v", err)
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
