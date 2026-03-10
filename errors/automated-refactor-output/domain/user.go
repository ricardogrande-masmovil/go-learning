package domain

// UserErrorCause represents the semantic type of user domain error.
type UserErrorCause int

const (
	UserErrorCauseUnknown UserErrorCause = iota
	UserErrorCauseNotFound
)

// UserDomainError represents a domain error strictly adhering to the
// split-brain DDD pattern.
type UserDomainError struct {
	UserMsg string
	Err     error
	Cause   UserErrorCause
}

// Error implements the error interface.
func (e *UserDomainError) Error() string {
	return e.UserMsg
}

// Unwrap allows checking for the underlying root cause if needed internally.
func (e *UserDomainError) Unwrap() error {
	return e.Err
}

type User struct {
	ID    string
	Email string
	Name  string
}
