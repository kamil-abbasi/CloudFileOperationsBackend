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

func (repository *FilesPostgresRepository) Find() ([]entities.File, error) {
	var files []entities.File

	rows, err := repository.db.Query("SELECT * FROM files")

	if err != nil {
		return nil, fmt.Errorf("Failed to find files: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var file entities.File

		err := rows.Scan(&file.Id, &file.Filename, &file.Location, &file.Size, &file.UserId)

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

	err := repository.db.QueryRow("SELECT * FROM files WHERE id = ?", id).Scan(&file.Id, &file.Filename, &file.Location, &file.Size, &file.UserId)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("Failed to find one file. Details: %v", err)
	}

	return file, true, nil
}

func (repository *FilesPostgresRepository) Create(dto dtos.FileCreateDto) (entities.File, error) {
	_, err := repository.db.Exec("INSERT INTO files VALUES(?,?,?,?,?)", dto.Id, dto.Filename, dto.Location, dto.Size, dto.UserId)

	if err != nil {
		return entities.File{}, fmt.Errorf("failed to create file metadata, details: %v", err)
	}

	return entities.File(dto), nil
}

func (repository *FilesPostgresRepository) Update(dto dtos.FileUpdateDto) (bool, error) {
	result, err := repository.db.Exec("UPDATE files SET filename = ?, location = ? WHERE id = ?", dto.Fields.Filename, dto.Fields.Location, dto.Where.Id)

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
	result, err := repository.db.Exec("DELETE FROM files WHERE id = ?", id)

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
