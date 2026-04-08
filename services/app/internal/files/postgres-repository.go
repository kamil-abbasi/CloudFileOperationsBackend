package files

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/pkg"
)

type FilesPostgresRepository struct {
	db *database.Postgres
}

func NewPostgresRepository(db *database.Postgres) IFilesRepository {
	return &FilesPostgresRepository{
		db: db,
	}
}

func (r *FilesPostgresRepository) Find(dto dtos.FindFilesDto) ([]entities.File, error) {
	directoryID, err := parseOptionalUUID(dto.Where.DirectoryId)

	if err != nil {
		return nil, fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	files, err := r.db.Queries.FindFiles(context.Background(), directoryID)

	if err != nil {
		return nil, fmt.Errorf("failed to find files: %v", err)
	}

	result := pkg.SliceMap(files, DatabaseToEntity)

	return result, nil
}

func (r *FilesPostgresRepository) FindByLocation(location string) ([]entities.File, error) {
	files, err := r.db.Queries.FindFilesByLocation(context.Background(), sql.NullString{String: location, Valid: true})

	if err != nil {
		return nil, fmt.Errorf("failed to find files: %v", err)
	}

	result := pkg.SliceMap(files, DatabaseToEntity)

	return result, nil
}

func (r *FilesPostgresRepository) FindOne(id string) (entities.File, bool, error) {
	if id == "" {
		return entities.File{}, false, nil
	}

	fileID, err := uuid.Parse(id)

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to parse file id, details: %v", err)
	}

	file, err := r.db.Queries.FindFileById(context.Background(), fileID)

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to find one file, details: %v", err)
	}

	return DatabaseToEntity(file), true, nil
}

func (r *FilesPostgresRepository) FindByNameAndDirectoryId(name string, directoryId string) (entities.File, bool, error) {
	queryDirID, err := parseOptionalUUID(directoryId)

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	file, err := r.db.Queries.FindFileByNameAndDirectoryId(context.Background(), database.FindFileByNameAndDirectoryIdParams{
		Name:        name,
		DirectoryID: queryDirID,
	})

	if err == sql.ErrNoRows {
		return entities.File{}, false, nil
	}

	if err != nil {
		return entities.File{}, false, fmt.Errorf("failed to find file by name and directory id, details: %v", err)
	}

	return DatabaseToEntity(file), true, nil
}

func (r *FilesPostgresRepository) Save(file entities.File) error {
	fileID, err := uuid.Parse(file.Id)

	if err != nil {
		return fmt.Errorf("failed to parse file id, details: %v", err)
	}

	userID, err := uuid.Parse(file.UserId)

	if err != nil {
		return fmt.Errorf("failed to parse user id, details: %v", err)
	}

	directoryID, err := parseOptionalUUID(file.DirectoryId)

	if err != nil {
		return fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	if file.Size > math.MaxInt32 {
		return fmt.Errorf("failed to save file to postgres, details: file size exceeds postgres integer max")
	}

	ctx := context.Background()

	tx, err := r.db.Conn.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	queries := r.db.Queries.WithTx(tx)

	err = queries.SaveFile(ctx, database.SaveFileParams{
		ID:          fileID,
		UserID:      userID,
		DirectoryID: directoryID,
		Name:        file.Name,
		Size:        int32(file.Size),
		Location:    file.Location,
		Checksum:    file.Checksum,
	})

	if err != nil {
		return fmt.Errorf("failed to save file to postgres, details: %v", err)
	}

	err = queries.SaveFileItem(ctx, database.SaveFileItemParams{
		ID:     fileID,
		UserID: userID,
	})

	if err != nil {
		return fmt.Errorf("failed to save file to postgres, details: %v", err)
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r *FilesPostgresRepository) Remove(id string) (bool, error) {
	fileID, err := uuid.Parse(id)

	if err != nil {
		return false, fmt.Errorf("failed to parse file id, details: %v", err)
	}

	ctx := context.Background()

	tx, err := r.db.Conn.BeginTx(ctx, nil)

	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	queries := r.db.Queries.WithTx(tx)

	rowsAffected, err := queries.RemoveFileRows(ctx, fileID)

	if err != nil {
		return false, fmt.Errorf("failed to remove file from postgres, details: %v", err)
	}

	err = queries.RemoveFileItem(ctx, fileID)

	if err != nil {
		return false, fmt.Errorf("failed to remove file from postgres, details: %v", err)
	}

	err = tx.Commit()

	if err != nil {
		return false, err
	}

	if rowsAffected <= 0 {
		return false, nil
	}

	return true, nil
}

func parseOptionalUUID(value string) (uuid.NullUUID, error) {
	if value == "" {
		return uuid.NullUUID{Valid: false}, nil
	}

	parsedUUID, err := uuid.Parse(value)

	if err != nil {
		return uuid.NullUUID{}, err
	}

	return uuid.NullUUID{UUID: parsedUUID, Valid: true}, nil
}
