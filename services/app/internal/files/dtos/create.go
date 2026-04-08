package dtos

type CreateFileDto struct {
	DirectoryId string `json:"directoryId"`
	Name        string `json:"name" binding:"required"`
	Checksum    string `json:"checksum" binding:"required"`
}
