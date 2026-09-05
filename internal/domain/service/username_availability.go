package service

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrUsernameAlreadyExists = errors.New(
		"username already exists",
	)
)

type UsernameAvailability struct {
	users repo.User
}

func NewUsernameAvailability(users repo.User) *UsernameAvailability {
	return &UsernameAvailability{
		users: users,
	}
}

func (s *UsernameAvailability) EnsureAvailable(
	ctx context.Context,
	username vo.Username,
) error {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	if user != nil {
		return ErrUsernameAlreadyExists
	}

	return nil
}
