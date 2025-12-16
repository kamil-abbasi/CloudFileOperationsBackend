package dtos

type FileUpdateDto struct {
	Where struct {
		Id string
	}
	Fields struct {
		Filename string
		Location string
	}
}
