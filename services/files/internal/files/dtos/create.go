package dtos

type FileCreateDto struct {
	Id       string
	UserId   string
	Filename string
	Location string
	Size     uint64
}
