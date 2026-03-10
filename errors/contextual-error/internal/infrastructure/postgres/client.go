package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain"
	"github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/internal/domain/model"
	ctxErr "github.com/ricardogrande-masmovil/go-learning/errors/contextual-error/pkg/contextual_errors"
)

// imagine these are real postgres errors
var (
	ErrConnectionRefused   = errors.New("connection refused")
	ErrPrimaryKeyViolation = errors.New("primary key violation")
)

type client struct {
	db *sql.DB
}

func NewClient(db *sql.DB) domain.UserRepository {
	return &client{db: db}
}

func (c *client) FindByID(ctx context.Context, id string) (*model.User, error) {

	return nil, ctxErr.Wrap(ErrConnectionRefused, "FindByID").Err()
}

func (c *client) Create(ctx context.Context, user *model.User) error {
	return ctxErr.Wrap(ErrPrimaryKeyViolation, "Create").Err()
}
