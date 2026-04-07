package files

import "errors"

var ErrNotFound = errors.New("file not found")
var ErrAlreadyExists = errors.New("file already exists")
