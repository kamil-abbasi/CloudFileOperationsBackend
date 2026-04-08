package dtos

import "mime/multipart"

type UploadFileDto struct {
	IdempotencyKey string                `form:"idempotencyKey" binding:"required"`
	DirectoryId    string                `form:"directoryId"`
	File           *multipart.FileHeader `form:"file" binding:"required"`
	Checksum       string                `form:"checksum" binding:"required"`
}
