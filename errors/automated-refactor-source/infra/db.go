package infra

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-source/domain"
)

// FindUser simulates a database lookup.
func FindUser(id string) (*domain.User, error) {
	// Simulate a database failure
	dbErr := sql.ErrNoRows

	if dbErr != nil {
		if errors.Is(dbErr, sql.ErrNoRows) {
			// NON-COMPLIANT: Opaque return. Returning a domain error without
			// wrapping the underlying root cause (`dbErr`), losing the error chain.
			return nil, domain.ErrUserNotFound
		}

		// NON-COMPLIANT: "failed to" anti-pattern and using %v which obfuscates
		// the error when we might want to preserve the chain here.
		return nil, fmt.Errorf("failed to query postgres database: %v", dbErr)
	}

	return &domain.User{ID: id, Email: "test@example.com"}, nil
}
