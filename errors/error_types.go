package errors

import "errors"

// Predefined error types
var (
	ErrBindingJSON     = errors.New("ERR_BINDING_JSON")
	ErrUserExists      = errors.New("ERR_USER_EXISTS")
	ErrUserNotFound    = errors.New("ERR_USER_NOT_FOUND")
	ErrInvalidPassword = errors.New("ERR_INVALID_PASSWORD")
	ErrHashingPassword = errors.New("ERR_HASHING_PASSWORD")
	ErrCreatingUser    = errors.New("ERR_CREATING_USER")
	ErrGeneratingToken = errors.New("ERR_GENERATING_TOKEN")
	ErrObjectExists    = errors.New("ERR_OBJECT_EXISTS")
)
