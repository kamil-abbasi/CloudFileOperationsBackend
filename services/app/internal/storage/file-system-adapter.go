package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type FileSystemAdapter struct {
	config *config.Config
}

func NewFileSystemAdapter(config *config.Config) IStorage {
	return &FileSystemAdapter{
		config: config,
	}
}

func (storage *FileSystemAdapter) FileExists(key string) (bool, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	return shared.DoesFileExist(path)
}

func (storage *FileSystemAdapter) UploadFile(key string, src io.Reader) (int64, string, error) {
	path := filepath.Clean(filepath.Join(storage.config.RootPath, key))

	exists, err := shared.DoesFileExist(path)

	if err != nil {
		return 0, "", err
	}

	if exists {
		return 0, "", ErrAlreadyExists
	}

	hash := sha256.New()
	tee := io.TeeReader(src, hash)

	file, err := os.Create(path)

	if err != nil {
		return 0, "", err
	}

	defer file.Close()

	bytesWritten, err := io.Copy(file, tee)

	if err != nil {
		return 0, "", err
	}

	finalHash := hex.EncodeToString(hash.Sum(nil))

	return bytesWritten, finalHash, nil
}

func (storage *FileSystemAdapter) DownloadFile(key string) (io.ReadCloser, bool, error) {
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

func (storage *FileSystemAdapter) RemoveFile(key string) (bool, error) {
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
