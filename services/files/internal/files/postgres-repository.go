package files

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type FilesPostgresRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewPostgresRepository(config *config.Config, db *sql.DB) interfaces.IFilesRepository {
	return &FilesPostgresRepository{
		db:     db,
		config: config,
	}
}

func (repository *FilesPostgresRepository) Find(dto dtos.FindFilesDto) ([]entities.File, error) {
	var files []entities.File

	rows, err := repository.db.Query(`
		SELECT id, user_id, directory_id, name, size
		FROM files
		WHERE directory_id = $1`, dto.Where.DirectoryId)

	if err != nil {
		return nil, fmt.Errorf("Failed to find files: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var file entities.File
		var directoryId sql.NullString

		err := rows.Scan(&file.Id, &file.UserId, &directoryId, &file.Name, &file.Size)

		if err != nil {
			return nil, err
		}

		if directoryId.Valid {
			file.DirectoryId = directoryId.String
		} else {
			file.DirectoryId = ""
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (repository *FilesPostgresRepository) FindOne(id string) (entities.File, bool, error) {
	if id == "" {
		return entities.File{}, false, nil
	}

	var file entities.File
	var directoryId sql.NullString

	err := repository.db.QueryRow(`
		SELECT id, user_id, directory_id, name, size
		FROM files
		WHERE id = $1`,
		id).Scan(&file.Id, &file.UserId, &directoryId, &file.Name, &file.Size)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("Failed to find one file. Details: %v", err)
	}

	if directoryId.Valid {
		file.DirectoryId = directoryId.String
	} else {
		file.DirectoryId = ""
	}

	return file, true, nil
}

func (repository *FilesPostgresRepository) FindByNameAndDirectoryId(name string, directoryId string) (entities.File, bool, error) {
	var file entities.File
	var dbDirectoryId sql.NullString

	var queryDirId sql.NullString

	if directoryId != "" {
		queryDirId = sql.NullString{
			String: directoryId,
			Valid:  true,
		}
	} else {
		queryDirId = sql.NullString{Valid: false}
	}

	err := repository.db.QueryRow(`
		SELECT id, user_id, directory_id, name, size
		FROM files
		WHERE name = $1 AND directory_id IS NOT DISTINCT FROM $2`, name, queryDirId).Scan(&file.Id, &file.UserId, &dbDirectoryId, &file.Name, &file.Size)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("Failed to find file by name and directory id. Details: %v", err)
	}

	if dbDirectoryId.Valid {
		file.DirectoryId = dbDirectoryId.String
	} else {
		file.DirectoryId = ""
	}

	return file, true, nil
}

func (repository *FilesPostgresRepository) Save(file entities.File) error {
	_, err := repository.db.Exec(
		`INSERT INTO
		files(id, user_id, directory_id, name, size)
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT(id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			directory_id = EXCLUDED.directory_id,
			name = EXCLUDED.name,
			size = EXCLUDED.size
		RETURNING id, user_id, directory_id, name, size`,
		file.Id, file.UserId, &sql.NullString{
			String: file.DirectoryId,
		}, file.Name, file.Size)

	if err != nil {
		return fmt.Errorf("failed to save file to postgres, details: %v", err)
	}

	return nil
}

func (repository *FilesPostgresRepository) Remove(id string) (bool, error) {
	result, err := repository.db.Exec("DELETE FROM files WHERE id = $1", id)

	if err != nil {
		return false, fmt.Errorf("failed to remove file from postgres, details: %v", err)
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
