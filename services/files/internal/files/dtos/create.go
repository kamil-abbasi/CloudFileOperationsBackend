package dtos

type CreateFileDto struct {
	Id          string
	UserId      string
	DirectoryId string
	Name        string
	Location    string
	Size        uint64
}
