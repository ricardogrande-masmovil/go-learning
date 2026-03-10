package main

import (
	"errors"
	"fmt"
)

// ErrNotFound is a distinct, static error value.
// Best used for sentinel errors that callers can check against using errors.Is.
var ErrNotFound = errors.New("not found")

// QueryError is a custom error type that provides structured "flavor".
// It can carry additional context like an error code or specific fields.
type QueryError struct {
	Message string
	Cause   ErrorCause
}

// Error implements the error interface for QueryError.
// This formats the final error message exposing the context.
func (e *QueryError) Error() string {
	return fmt.Sprintf("query error (cause %d): %s", e.Cause, e.Message)
}

// ErrorCause represents specific reasons for a query failure.
type ErrorCause int

const (
	NetworkError ErrorCause = iota + 1
	UserError
	UnknownError
)

// FetchUser simulates fetching a user from a database.
// It demonstrates different ways of creating and returning errors.
func FetchUser(id int) error {
	if id == 0 {
		// Using errors.New for a simple, distinct sentinel error.
		return ErrNotFound
	}

	if id < 0 {
		// Using fmt.Errorf for a dynamic, formatted string.
		return fmt.Errorf("user id %d is invalid: must be positive", id)
	}

	// Returning a custom error type for structured details and flavor.
	return &QueryError{
		Message: "database connection lost",
		Cause:   NetworkError, // 1
	}
}

func main() {
	fmt.Println("=== 01: Raw Ingredients of a Go Error ===")

	err1 := FetchUser(0)
	fmt.Printf("Sentinel Error: %v\n", err1)

	err2 := FetchUser(-5)
	fmt.Printf("Formatted Error: %v\n", err2)

	err3 := FetchUser(42)
	fmt.Printf("Custom Error Type: %v\n", err3)
}
