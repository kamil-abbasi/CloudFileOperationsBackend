package directories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/pkg"
)

type DirectoriesPostgresRepository struct {
	db *database.Postgres
}

func NewPostgresRepository(db *database.Postgres) IDirectoriesRepository {
	return &DirectoriesPostgresRepository{
		db: db,
	}
}

func (r *DirectoriesPostgresRepository) ListItems(location string) ([]entities.DirectoryItem, error) {
	items, err := r.db.Queries.ListDirectoryItems(context.Background(), location)

	if err != nil {
		return nil, fmt.Errorf("failed to find directory items: %v", err)
	}

	result := pkg.SliceMap(items, DatabaseItemToEntity)

	return result, nil
}

func (r *DirectoriesPostgresRepository) Find() ([]entities.Directory, error) {
	directories, err := r.db.Queries.FindDirectories(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to find directories: %v", err)
	}

	result := pkg.SliceMap(directories, DatabaseToEntity)

	return result, nil
}

func (r *DirectoriesPostgresRepository) FindOne(id string) (entities.Directory, bool, error) {
	if id == "" {
		return entities.Directory{}, false, nil
	}

	directoryID, err := uuid.Parse(id)

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	directory, err := r.db.Queries.FindDirectoryById(context.Background(), directoryID)

	if err == sql.ErrNoRows {
		return entities.Directory{}, false, nil
	}

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("failed to find one directory, details: %v", err)
	}

	return DatabaseToEntity(directory), true, nil
}

func (r *DirectoriesPostgresRepository) FindByNameAndParentId(name string, parentId string) (entities.Directory, bool, error) {
	queryParentID, err := parseOptionalUUID(parentId)

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("failed to parse parent id, details: %v", err)
	}

	directory, err := r.db.Queries.FindDirectoryByNameAndParentId(context.Background(), database.FindDirectoryByNameAndParentIdParams{
		Name:     name,
		ParentID: queryParentID,
	})

	if err == sql.ErrNoRows {
		return entities.Directory{}, false, nil
	}

	if err != nil {
		return entities.Directory{}, false, fmt.Errorf("failed to find directory by name and parent id, details: %v", err)
	}

	return DatabaseToEntity(directory), true, nil
}

func (r *DirectoriesPostgresRepository) Save(directory entities.Directory) error {
	directoryID, err := uuid.Parse(directory.Id)

	if err != nil {
		return fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	userID, err := uuid.Parse(directory.UserId)

	if err != nil {
		return fmt.Errorf("failed to parse user id, details: %v", err)
	}

	parentID, err := parseOptionalUUID(directory.ParentId)

	if err != nil {
		return fmt.Errorf("failed to parse parent id, details: %v", err)
	}

	ctx := context.Background()

	tx, err := r.db.Conn.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	queries := r.db.Queries.WithTx(tx)

	err = queries.SaveDirectory(ctx, database.SaveDirectoryParams{
		ID:       directoryID,
		UserID:   userID,
		ParentID: parentID,
		Name:     directory.Name,
		Location: directory.Location,
	})

	if err != nil {
		return fmt.Errorf("failed to save directory to postgres, details: %v", err)
	}

	err = queries.SaveDirectoryItem(ctx, database.SaveDirectoryItemParams{
		ID:     directoryID,
		UserID: userID,
	})

	if err != nil {
		return fmt.Errorf("failed to save directory to postgres, details: %v", err)
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r *DirectoriesPostgresRepository) Remove(id string) (bool, error) {
	directoryID, err := uuid.Parse(id)

	if err != nil {
		return false, fmt.Errorf("failed to parse directory id, details: %v", err)
	}

	ctx := context.Background()

	tx, err := r.db.Conn.BeginTx(ctx, nil)

	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	queries := r.db.Queries.WithTx(tx)

	rowsAffected, err := queries.RemoveDirectoryRows(ctx, directoryID)

	if err != nil {
		return false, fmt.Errorf("failed to remove directory from postgres, details: %v", err)
	}

	err = queries.RemoveDirectoryItem(ctx, directoryID)

	if err != nil {
		return false, fmt.Errorf("failed to remove directory from postgres, details: %v", err)
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
