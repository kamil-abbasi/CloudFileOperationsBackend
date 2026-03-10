package dtos

type CreateDirectoryDto struct {
	Id       string
	UserId   string
	ParentId string
	Location string
	Name     string
}
