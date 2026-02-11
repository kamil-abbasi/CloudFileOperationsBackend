package files

import (
	"database/sql"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/interfaces"
)

type FilesPostgresRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewPostgresRepository(config *config.Config, db *sql.DB) interfaces.IFilesRepository {
	return &FilesSQLiteRepository{
		db:     db,
		config: config,
	}
}

func (repository *FilesPostgresRepository) Find() ([]entities.File, error) {
	return nil, nil
}

func (repository *FilesPostgresRepository) FindOne(id string) (*entities.File, error) {
	return nil, nil
}

func (repository *FilesPostgresRepository) Create() (*entities.File, error) {
	return nil, nil
}

func (repository *FilesPostgresRepository) Update() (*entities.File, error) {
	return nil, nil
}

func (repository *FilesPostgresRepository) Remove(id string) (bool, error) {
	return false, nil
}
