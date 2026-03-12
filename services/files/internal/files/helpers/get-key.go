package helpers

import (
	"fmt"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files/entities"
)

func GetFileKey(file entities.File) string {
	return fmt.Sprintf("%s-%s-%s", file.UserId, file.DirectoryId, file.Id)
}
