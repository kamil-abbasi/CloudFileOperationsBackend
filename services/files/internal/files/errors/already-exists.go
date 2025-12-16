package errors

type FileAlreadyExistsError struct {
	Err error
}

func (e *FileAlreadyExistsError) Error() string {
	return e.Err.Error()
}
