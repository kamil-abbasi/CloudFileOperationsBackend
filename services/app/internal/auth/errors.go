package auth

import "errors"

var ErrPasswordsDoNotMatch = errors.New("passwords do not match")
var ErrUserAlreadyExists = errors.New("user with this name already exists")
var ErrHashGeneration = errors.New("failed to generate password hash")
var ErrUserNotFound = errors.New("user with this name does not exist")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrMissingAuthHeader = errors.New("missing auth header")
var ErrMissingToken = errors.New("missing token")
var ErrInvalidToken = errors.New("invalid token")
