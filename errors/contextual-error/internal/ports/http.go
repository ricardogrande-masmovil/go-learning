package ports

import (
	"context"

	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain/model"
	ctxErr "github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/pkg/contextual_errors"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, id, name string) error
}

// Fake http handler to demonstrate error handling
type HttpHandler struct {
	userService UserService
}

func NewHTTPHandler(userService UserService) *HttpHandler {
	return &HttpHandler{userService: userService}
}

// GetUserByID is a fake http handler to demonstrate error handling, it returns an http response
func (h *HttpHandler) GetUserByID(ctx context.Context, id string) {
	_, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		log.Error().
			Err(err).
			Fields(ctxErr.GetMetadata(err)).
			Msg("failed to get user by ID in http handler")
		return
	}
	log.Info().Msg("HTTP 200: User retrieved successfully")
}

// CreateUser is a fake http handler to demonstrate error handling, it returns an http response
func (h *HttpHandler) CreateUser(ctx context.Context, id, name string) {
	err := h.userService.CreateUser(ctx, id, name)
	if err != nil {
		log.Error().
			Err(err).
			Fields(ctxErr.GetMetadata(err)).
			Msg("failed to create user in http handler")
		return
	}
	log.Info().Msg("HTTP 201: User created successfully")
}
