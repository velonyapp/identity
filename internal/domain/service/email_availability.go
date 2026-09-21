package service

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrEmailAlreadyExists = errors.New(
		"email already exists",
	)
)

type EmailAvailability struct {
	users repo.User
}

func NewEmailAvailability(users repo.User) *EmailAvailability {
	return &EmailAvailability{
		users: users,
	}
}

func (s *EmailAvailability) EnsureAvailable(
	ctx context.Context,
	email vo.Email,
) error {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user != nil {
		return ErrEmailAlreadyExists
	}

	return nil
}
