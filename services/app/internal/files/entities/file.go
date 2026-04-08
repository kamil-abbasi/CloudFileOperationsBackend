package entities

type File struct {
	Id          string
	UserId      string
	DirectoryId string
	Name        string
	Size        uint64
	Location    string
	Checksum    string
}
