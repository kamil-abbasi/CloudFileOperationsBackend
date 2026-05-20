package files

import "errors"

var ErrNotFound = errors.New("file not found")
var ErrAlreadyExists = errors.New("file already exists")
var ErrDirectoryNotFound = errors.New("directory not found")
var ErrCorruptedUpload = errors.New("corrupted upload")
var ErrUserNotFound = errors.New("user not found")
var ErrNotEnoughStorage = errors.New("not enough storage")
