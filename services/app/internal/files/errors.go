package files

import "errors"

var ErrNotFound = errors.New("file not found")
var ErrAlreadyExists = errors.New("file already exists")
var ErrDirectoryNotFound = errors.New("directory not found")
