package auth

import (
	"errors"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"

	"golang.org/x/crypto/bcrypt"
)

var _ port.PasswordManager = (*passwordManager)(nil)

type passwordManager struct {
	cost int
}

func NewPasswordManager() port.PasswordManager {
	return &passwordManager{cost: bcrypt.DefaultCost}
}

func (h *passwordManager) Hash(password vo.Password) (vo.PasswordHash, error) {
	rawPasswordHash, err := bcrypt.GenerateFromPassword([]byte(password.String()), h.cost)
	if err != nil {
		return vo.PasswordHash{}, err
	}

	passwordHash, err := vo.NewPasswordHash(rawPasswordHash)
	if err != nil {
		return vo.PasswordHash{}, err
	}

	return vo.PasswordHash(passwordHash), nil
}

func (h *passwordManager) Verify(password vo.Password, passwordHash vo.PasswordHash) error {
	err := bcrypt.CompareHashAndPassword(passwordHash.Bytes(), []byte(password.String()))
	if err == nil {
		return nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return port.ErrPasswordMismatch
	}

	return err
}
