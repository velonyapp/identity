package vo

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrInvalidEmail = errors.New(
		"email is invalid",
	)
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	addr, err := mail.ParseAddress(value)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}

	if addr.Address != value {
		return Email{}, ErrInvalidEmail
	}

	return Email{
		value: value,
	}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}
