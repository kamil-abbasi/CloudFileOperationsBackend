package dtos

type FileUpdateDto struct {
	Id     string
	Fields struct {
		Filename string
		Location string
	}
}
