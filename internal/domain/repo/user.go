package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserVersionConflict   = errors.New("user version conflict")
)

type User interface {
	FindByID(ctx context.Context, userID vo.UserID) (*entity.User, error)
	FindByUsername(ctx context.Context, username vo.Username) (*entity.User, error)
	FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error)

	Save(ctx context.Context, user *entity.User) error
}
