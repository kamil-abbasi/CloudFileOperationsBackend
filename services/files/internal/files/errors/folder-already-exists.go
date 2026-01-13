package errors

type FolderAlreadyExistsError struct {
	Err error
}

func (e *FolderAlreadyExistsError) Error() string {
	if e.Err == nil {
		return "folder already exists"
	}
	return e.Err.Error()
}
