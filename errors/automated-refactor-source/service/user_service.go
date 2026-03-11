package service

import (
	"fmt"
	"log"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-source/domain"
	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-source/infra"
)

// GetUserProfile simulates a domain service layer coordinating data retrieval.
func GetUserProfile(id string) (*domain.User, error) {
	user, err := infra.FindUser(id)
	if err != nil {
		// NON-COMPLIANT: Double-logging. We log the error here using the standard logger...
		log.Printf("[ERROR] issue fetching user %s: %v", id, err)

		// ...and then we return it again up the stack.
		// NON-COMPLIANT: The "failed to" string wrapping anti-pattern is also present.
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return user, nil
}
