package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/velony-app/identity/internal/application/port"
	"github.com/velony-app/identity/internal/domain/vo"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher() port.PasswordHasher {
	return &PasswordHasher{cost: bcrypt.DefaultCost}
}

func (h *PasswordHasher) Hash(password vo.Password) (vo.PasswordHash, error) {
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

func (h *PasswordHasher) Verify(password vo.Password, passwordHash vo.PasswordHash) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash.Value()), []byte(password.Value()))
	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return port.ErrPasswordMismatch
	}

	return err
}
