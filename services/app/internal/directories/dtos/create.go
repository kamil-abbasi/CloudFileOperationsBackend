package dtos

type CreateDirectoryDto struct {
	ParentId string `json:"parentId"`
	Name     string `json:"name" binding:"required"`
}
