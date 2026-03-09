package infra

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated_refactor_output/domain"
)

// FindUser simulates a database lookup.
func FindUser(id string) (*domain.User, error) {
	// Simulate a database failure
	dbErr := sql.ErrNoRows

	if dbErr != nil {
		cause := domain.UserErrorCauseUnknown
		userMsg := "An unexpected error occurred while looking up the user."

		if errors.Is(dbErr, sql.ErrNoRows) {
			cause = domain.UserErrorCauseNotFound
			userMsg = "User could not be found."
		}

		return nil, &domain.UserDomainError{
			UserMsg: userMsg,
			Err:     fmt.Errorf("query postgres database: %w", dbErr),
			Cause:   cause,
		}
	}

	return &domain.User{ID: id, Email: "test@example.com"}, nil
}
