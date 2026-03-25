package entities

type ItemType string

const (
	File ItemType = "file"
	Dir  ItemType = "directory"
)

type DirectoryItem struct {
	Id     string
	UserId string
	Type   ItemType
	Name   string
}
