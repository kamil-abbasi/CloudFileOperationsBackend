package storage

import "io"

type IStorage interface {
	FileExists(key string) (bool, error)
	UploadFile(key string, src io.Reader) (int64, string, error)
	RemoveFile(key string) (bool, error)
	DownloadFile(key string) (io.ReadCloser, bool, error)
}
