package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FileSystemStorageAdapter struct {
	config *config.Config
}

func NewFileSystemAdapter(config *config.Config) IStorage {
	return &FileSystemStorageAdapter{
		config: config,
	}
}

func (storage *FileSystemStorageAdapter) FileExists(key string) (bool, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	return shared.DoesFileExist(path)
}

func (storage *FileSystemStorageAdapter) UploadFile(key string, src io.Reader) (int64, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	exists, err := shared.DoesFileExist(path)

	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrAlreadyExists
	}

	file, err := os.Create(path)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	bytesWritten, err := io.Copy(file, src)

	if err != nil {
		return 0, err
	}

	return bytesWritten, nil
}

func (storage *FileSystemStorageAdapter) DownloadFile(key string) (io.ReadCloser, bool, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	exists, err := shared.DoesFileExist(path)

	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, false, nil
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, false, err
	}

	return file, true, nil
}

func (storage *FileSystemStorageAdapter) RemoveFile(key string) (bool, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	err := os.Remove(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
