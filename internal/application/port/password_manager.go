package port

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrPasswordMismatch = errors.New("password does not match")

type PasswordManager interface {
	Hash(ctx context.Context, password vo.Password) (vo.PasswordHash, error)
	Verify(ctx context.Context, password vo.Password, hash vo.PasswordHash) error
}
