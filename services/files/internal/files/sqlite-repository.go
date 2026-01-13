package files

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type FilesSQLiteRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewSQLiteRepository(config *config.Config) interfaces.IFilesRepository {
	db := database.GetInstance()

	return &FilesSQLiteRepository{
		db:     db,
		config: config,
	}
}
//abcd
func (repository *FilesSQLiteRepository) FindOne(id string) (entities.File, bool, error) {
	var file entities.File

	err := repository.db.QueryRow(
		"SELECT id, filename, location, size, user_id FROM files WHERE id = ?",
		id,
	).Scan(&file.Id, &file.Filename, &file.Location, &file.Size, &file.UserId)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to find file, details: %v", err)
	}

	return file, true, nil
}
func (repository *FilesSQLiteRepository) Create(dto dtos.FileCreateDto) (entities.File, error) {

	_, err := repository.db.Exec("INSERT INTO files VALUES(?,?,?,?,?)", dto.Id, dto.Filename, dto.Location, dto.Size, dto.UserId)

	if err != nil {
		return entities.File{}, fmt.Errorf("failed to create file metadata, details: %v", err)
	}

	return entities.File(dto), nil
}

func (repository *FilesSQLiteRepository) CreateFolder(dto dtos.FolderCreateDto) (entities.Folder, error) {
	_, err := repository.db.Exec("INSERT INTO folders VALUES(?,?,?,?)", dto.Id, dto.Name, dto.Location, dto.UserId)

	if err != nil {
		return entities.Folder{}, fmt.Errorf("failed to create folder metadata, details: %v", err)
	}

	return entities.Folder(dto), nil
}
func (r *FilesSQLiteRepository) Update(dto dtos.FileUpdateDto) (bool, error) {
	result, err := r.db.Exec("UPDATE files SET filename = ?, location = ? WHERE id = ?", dto.Fields.Filename, dto.Fields.Location, dto.Where.Id)

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
func (repository *FilesSQLiteRepository) Remove(id string) (bool, error) {
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
