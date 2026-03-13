package shared

type HttpError struct {
	Message string
	Code    uint16
	Details map[string]any
}

type FileAlreadyExistsError struct {
	Err error
}

func (e *FileAlreadyExistsError) Error() string {
	return e.Err.Error()
}

type DirectoryNotFoundError struct {
	Err error
}

func (e *DirectoryNotFoundError) Error() string {
	return e.Err.Error()
}

type DirectoryAlreadyExistsError struct {
	Err error
}

func (e *DirectoryAlreadyExistsError) Error() string {
	return e.Err.Error()
}
