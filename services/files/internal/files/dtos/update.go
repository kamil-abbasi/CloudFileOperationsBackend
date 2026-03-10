package dtos

type UpdateFileDto struct {
	Where struct {
		Id string
	}
	Fields struct {
		DirectoryId string
		Name        string
		Location    string
	}
}
