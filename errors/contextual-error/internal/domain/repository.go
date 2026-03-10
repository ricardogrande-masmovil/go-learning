package domain

import (
	"context"

	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}
