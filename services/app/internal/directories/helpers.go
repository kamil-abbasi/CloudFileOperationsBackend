package directories

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories/entities"
)

func EntityToDto(entity entities.Directory) dtos.DirectoryResponseDto {
	return dtos.DirectoryResponseDto{
		Id:       entity.Id,
		UserId:   entity.UserId,
		ParentId: entity.ParentId,
		Name:     entity.Name,
		Location: entity.Location,
	}
}

func ItemEntityToDto(entity entities.DirectoryItem) dtos.DirectoryItemResponseDto {
	return dtos.DirectoryItemResponseDto{
		Id:     entity.Id,
		UserId: entity.UserId,
		Name:   entity.Name,
		Type:   string(entity.Type),
	}
}

func DatabaseToEntity(directory database.Directory) entities.Directory {
	entity := entities.Directory{
		Id:       directory.ID.String(),
		UserId:   directory.UserID.String(),
		Name:     directory.Name,
		Location: directory.Location,
	}

	if directory.ParentID.Valid {
		entity.ParentId = directory.ParentID.UUID.String()
	}

	return entity
}

func DatabaseItemToEntity(item database.ListDirectoryItemsRow) entities.DirectoryItem {
	return entities.DirectoryItem{
		Id:     item.ID.String(),
		UserId: item.UserID.String(),
		Type:   entities.ItemType(item.Type),
		Name:   item.Name,
	}
}
