package storage

type IStorageAdapter interface {
	Upload()
	Remove()
	Download()
}
