package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.PasswordHasher = (*passwordHasher)(nil)

type passwordHasher struct {
	cost int
}

func NewPasswordHasher() port.PasswordHasher {
	return &passwordHasher{cost: bcrypt.DefaultCost}
}

func (h *passwordHasher) Hash(password vo.Password) (vo.PasswordHash, error) {
	rawPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password.Value()), h.cost)
	if err != nil {
		return vo.PasswordHash{}, err
	}

	passwordHash, err := vo.NewPasswordHash(string(rawPasswordHash))
	if err != nil {
		return vo.PasswordHash{}, err
	}

	return vo.PasswordHash(passwordHash), nil
}

func (h *passwordHasher) Verify(password vo.Password, passwordHash vo.PasswordHash) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash.Value()), []byte(password.Value()))
	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return port.ErrPasswordMismatch
	}

	return err
}
