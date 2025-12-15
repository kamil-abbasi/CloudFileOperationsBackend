package files

import (
	"database/sql"
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type FilesSQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() interfaces.IFilesRepository {
	db := database.GetInstance()

	return &FilesSQLiteRepository{
		db: db,
	}
}

func (repository *FilesSQLiteRepository) FindOne(id string) (entities.File, bool, error) {
	var file entities.File

	err := repository.db.QueryRow(
		"SELECT id, filename, location, size FROM files WHERE id = ?",
		id,
	).Scan(&file.Id, &file.Filename, &file.Location, &file.Size)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to find file, details: %v", err)
	}

	return file, true, nil
}
func (repository *FilesSQLiteRepository) Create(dto dtos.FileCreateDto) (entities.File, error) {

	_, err := repository.db.Exec("INSERT INTO files VALUES(?,?,?,?)", dto.Id, dto.Filename, dto.Location, dto.Size)

	if err != nil {
		return entities.File{}, fmt.Errorf("failed to create file metadata, details: %v", err)
	}

	return entities.File(dto), nil
}
func (repository *FilesSQLiteRepository) Update(dto dtos.FileUpdateDto) (entities.File, bool, error) {
	return entities.File{}, true, nil
}
func (repository *FilesSQLiteRepository) Remove(id string) (entities.File, bool, error) {
	return entities.File{}, true, nil
}
