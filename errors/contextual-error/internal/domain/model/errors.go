package model

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type ValidationCause string

const (
	CauseInvalidID   ValidationCause = "invalid_id"
	CauseInvalidName ValidationCause = "invalid_name"
)

// ValidationError represents a business rule validation failure.
// It does not support error chaining (i.e., it does not implement Unwrap() error)
// because validation errors are typically the root cause of an operation failing
// at the domain level, rather than a wrapper around another underlying error.
type ValidationError struct {
	Cause   ValidationCause
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
