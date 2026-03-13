package dtos

type UpdateFileDto struct {
	Where struct {
		Id string
	}
	Fields struct {
		Name        string
		DirectoryId string
	}
}
