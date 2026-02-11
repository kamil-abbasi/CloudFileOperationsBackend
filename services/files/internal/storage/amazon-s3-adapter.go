package storage

type AmazonS3StorageAdapter struct{}

func NewAmazonS3StorageAdapter() IStorageAdapter {
	return &AmazonS3StorageAdapter{}
}

func (storage AmazonS3StorageAdapter) Upload() {}

func (storage AmazonS3StorageAdapter) Download() {}

func (storage AmazonS3StorageAdapter) Remove() {}
