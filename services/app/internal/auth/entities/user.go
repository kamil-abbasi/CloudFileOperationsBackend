package entities

type User struct {
	Id           string
	Name         string
	PasswordHash string
	MaxStorage   uint64
	StorageUsed  uint64
}
