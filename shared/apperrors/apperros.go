package apperrors

import "errors"

// Common errors shared across all services
var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrInvalidRequest = errors.New("invalid request")
	ErrForbidden      = errors.New("forbidden")
	ErrInternal       = errors.New("internal error")
)
