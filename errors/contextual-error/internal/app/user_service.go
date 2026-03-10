package app

import (
	"context"
	"errors"

	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain"
	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain/model"
	ctxErr "github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/pkg/contextual_errors"
)

type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, ctxErr.Wrap(err, "GetUserByID").Str("user_id", id).Err()
	}
	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, id, name string) error {
	user, err := model.NewUser(id, name)
	if err != nil {
		if valErr, ok := errors.AsType[*model.ValidationError](err); ok {
			return ctxErr.Wrap(err, "CreateUser").Str("validation_cause", string(valErr.Cause)).Err()
		}
		return ctxErr.Wrap(err, "CreateUser").Err()
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return ctxErr.Wrap(err, "CreateUser").Str("user_id", user.ID()).Err()
	}
	return nil
}
