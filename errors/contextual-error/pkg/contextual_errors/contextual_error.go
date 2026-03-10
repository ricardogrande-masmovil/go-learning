package contextual_errors

import (
	"errors"
	"fmt"
)

// contextualError is an error that contains additional context.
type contextualError struct {
	err      error
	metadata map[string]any
}

func (e *contextualError) Error() string {
	return e.err.Error()
}

func (e *contextualError) Unwrap() error {
	return e.err
}

func (e *contextualError) Metadata() map[string]any {
	return e.metadata
}

// With creates a new contextual error with the given error.
func With(err error) *contextualError {
	return &contextualError{
		err:      err,
		metadata: make(map[string]any),
	}
}

// Wrap creates a new contextual error with the given error and message, wrapping the original error.
func Wrap(err error, msg string) *contextualError {
	if err == nil {
		return &contextualError{
			err:      nil,
			metadata: make(map[string]any),
		}
	}
	return &contextualError{
		err:      fmt.Errorf("%s: %w", msg, err),
		metadata: make(map[string]any),
	}
}

// Str adds a string metadata to the error.
func (e *contextualError) Str(key, value string) *contextualError {
	e.metadata[key] = value
	return e
}

// Int adds an int metadata to the error.
func (e *contextualError) Int(key string, value int) *contextualError {
	e.metadata[key] = value
	return e
}

// Float adds a float metadata to the error.
func (e *contextualError) Float(key string, value float64) *contextualError {
	e.metadata[key] = value
	return e
}

// Bool adds a bool metadata to the error.
func (e *contextualError) Bool(key string, value bool) *contextualError {
	e.metadata[key] = value
	return e
}

// Any adds an any metadata to the error.
func (e *contextualError) Any(key string, value any) *contextualError {
	e.metadata[key] = value
	return e
}

// Err returns the contextual error as an error.
// Useful to finalize the contextual error chaining API while avoiding the interface trap.
func (e *contextualError) Err() error {
	if e == nil || e.err == nil {
		return nil
	}
	return e
}

func GetMetadata(err error) map[string]any {
	data := make(map[string]any)

	for err != nil {
		if e, ok := err.(*contextualError); ok {
			for k, v := range e.metadata {
				data[k] = v
			}
		}

		err = errors.Unwrap(err)
	}

	return data
}
