package domain

import "errors"

// NON-COMPLIANT: This is a sentinel error defined in a domain layer.
// According to our DDD guidelines, this should be converted to a custom
// error struct (e.g., UserDomainError) that holds semantic information
// like the specific User ID and a Cause enum, rather than a flat string.
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    string
	Email string
	Name  string
}
