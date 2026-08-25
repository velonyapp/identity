package port

import (
	"errors"

	"github.com/velony-app/identity/internal/domain/vo"
)

var ErrPasswordMismatch = errors.New("password does not match")

type PasswordHasher interface {
	Hash(password vo.Password) (vo.PasswordHash, error)
	Verify(password vo.Password, hash vo.PasswordHash) error
}
