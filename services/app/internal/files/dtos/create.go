package dtos

type CreateFileDto struct {
	DirectoryId string `json:"directoryId"`
	Name        string `json:"name" binding:"required"`
}
