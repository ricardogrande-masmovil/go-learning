package main

import (
	"errors"
	"fmt"
)

// --- Domain / Public API Layer ---

// SafeError implements the Split-Brain pattern.
// Stop the infrastructure error chain at the repository boundary.
// Sanitize data before it crosses into the public dining room.
type SafeError struct {
	Code        string // Application specific code (e.g. "USER_NOT_FOUND")
	UserMsg     string // Sanitized domain error served for public exposure (like API responses)
	InternalErr error  // Stack trace used for logs
}

// Error implements the error interface.
func (e *SafeError) Error() string {
	return e.UserMsg
}

// --- Infrastructure Layer ---

// errPGAuthFailed simulates a raw database error containing secure details.
var errPGAuthFailed = errors.New("FATAL: password authentication failed for user 'admin'")

// DBRepository simulates our infrastructure layer accessing data.
type DBRepository struct{}

// FindUser simulates a DB access that fails.
func (repo *DBRepository) FindUser(id string) (*UserDTO, error) {
	// Raw DB or filesystem errors can leak SQL paths, library versions or tokens.
	return nil, errPGAuthFailed
}

// UserDTO is an infrastructure model based on DB schema.
// "And remember, errors are values! Would you return a DB DTO in the domain repository?"
// Context: Just like we don't leak DTOs, we shouldn't leak infra errors.
type UserDTO struct {
	ID   string
	Name string
}

// --- Service / Domain Layer ---

type UserService struct {
	repo *DBRepository
}

// GetUserProfile acts as the domain boundary.
func (s *UserService) GetUserProfile(id string) error {
	_, err := s.repo.FindUser(id)
	if err != nil {
		// Stop the infrastructure error chain at the repository boundary!
		// We do NOT return the raw error or simply wrap it with %w to the caller.
		// We map it to our SafeError domain model.

		return &SafeError{
			Code:        "INTERNAL_ERROR",
			UserMsg:     "We are unable to process your request at this time.",      // Sanitized
			InternalErr: fmt.Errorf("db failure when finding user %s: %w", id, err), // Internal logs only
		}
	}
	return nil
}

func main() {
	fmt.Println("=== 06: The Split-Brain Pattern (Domain Boundaries) ===")

	service := &UserService{repo: &DBRepository{}}

	// Make the call. Returns an error from the boundary.
	err := service.GetUserProfile("123")

	// At the HTTP Handler / Transport Layer, we receive the error:
	if err != nil {
		// We use type assertion / errors.As to check if this is our Safe domain error.
		var safeErr *SafeError
		if errors.As(err, &safeErr) {

			fmt.Println("--- Separation of Concerns ---")
			// 1. Used for backend logging (has the stack trace / raw db details)
			fmt.Printf("[LOG TO SYSTEM]             : %v\n", safeErr.InternalErr)

			// 2. Used for public API / JSON response (Sanitized)
			fmt.Printf("[HTTP RESPONSE to END USER] : HTTP 500 - {\"code\": \"%s\", \"message\": \"%s\"}\n",
				safeErr.Code, safeErr.UserMsg)
		} else {
			// Fallback for unknown/unhandled errors
			fmt.Println("[HTTP RESPONSE]: HTTP 500 - Unknown error")
		}
	}
}
