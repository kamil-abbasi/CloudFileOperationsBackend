package directories

import "errors"

var ErrNotFound = errors.New("directory not found")
var ErrAlreadyExists = errors.New("directory already exists")
