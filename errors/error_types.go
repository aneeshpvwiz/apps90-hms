package errors

import "errors"

// Predefined error types
var (
	ErrBindingJSON = errors.New("ERR_BINDING_JSON")
	ErrUserExists  = errors.New("ERR_USER_EXISTS")
)
