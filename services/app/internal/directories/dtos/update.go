package dtos

type UpdateDirectoryDto struct {
	Where struct {
		Id string
	}
	Fields struct {
		Name string
	}
}
