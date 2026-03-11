package service

import (
	"fmt"

	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-output/domain"
	"github.com/ricardogrande-masmovil/go-learning/errors/automated-refactor-output/infra"
)

// GetUserProfile simulates a domain service layer coordinating data retrieval.
func GetUserProfile(id string) (*domain.User, error) {
	user, err := infra.FindUser(id)
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}

	return user, nil
}
