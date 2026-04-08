package files

import (
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/database"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/dtos"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

func EntityToDto(entity entities.File) dtos.FileResponseDto {
	return dtos.FileResponseDto{
		Id:          entity.Id,
		UserId:      entity.UserId,
		DirectoryId: entity.DirectoryId,
		Name:        entity.Name,
		Size:        entity.Size,
		Location:    entity.Location,
		Checksum:    entity.Checksum,
	}
}

func DatabaseToEntity(file database.File) entities.File {
	entity := entities.File{
		Id:       file.ID.String(),
		UserId:   file.UserID.String(),
		Name:     file.Name,
		Size:     uint64(file.Size),
		Location: file.Location,
		Checksum: file.Checksum,
	}

	if file.DirectoryID.Valid {
		entity.DirectoryId = file.DirectoryID.UUID.String()
	}

	return entity
}
