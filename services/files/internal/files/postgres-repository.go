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

	rows, err := repository.db.Query("SELECT * FROM files WHERE directory_id = $1", dto.Where.DirectoryId)

	if err != nil {
		return nil, fmt.Errorf("Failed to find files: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var file entities.File

		err := rows.Scan(&file.Id, &file.Name, &file.Location, &file.Size, &file.UserId, &file.DirectoryId)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (repository *FilesPostgresRepository) FindOne(id string) (entities.File, bool, error) {
	var file entities.File

	err := repository.db.QueryRow("SELECT * FROM files WHERE id = $1", id).Scan(&file.Id, &file.Name, &file.Location, &file.Size, &file.UserId, &file.DirectoryId)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("Failed to find one file. Details: %v", err)
	}

	return file, true, nil
}

func (repository *FilesPostgresRepository) Create(dto dtos.CreateFileDto) (entities.File, error) {
	_, err := repository.db.Exec(
		`INSERT INTO
		files(id, name, location, size, user_id, directory_id)
		VALUES($1, $2, $3, $4, $5, $6)`,
		dto.Id, dto.Name, dto.Location, dto.Size, dto.UserId, dto.DirectoryId)

	if err != nil {
		return entities.File{}, fmt.Errorf("failed to create file metadata, details: %v", err)
	}

	return entities.File(dto), nil
}

func (repository *FilesPostgresRepository) Update(dto dtos.UpdateFileDto) (bool, error) {
	result, err := repository.db.Exec("UPDATE files SET name = $1, location = $2,directory_id = $3 WHERE id = $4", dto.Fields.Name, dto.Fields.Location, dto.Fields.DirectoryId, dto.Where.Id)

	if err != nil {
		return false, fmt.Errorf("failed to remove file metadata, details: %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, fmt.Errorf("failed to fetch the number of updated rows: %v", err)
	}

	if rowsAffected <= 0 {
		return false, nil
	}

	return true, nil
}

func (repository *FilesPostgresRepository) Remove(id string) (bool, error) {
	result, err := repository.db.Exec("DELETE FROM files WHERE id = $1", id)

	if err != nil {
		return false, fmt.Errorf("failed to remove file metadata, details: %v", err)
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
