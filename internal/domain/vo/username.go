package vo

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrUsernameTooShort = errors.New(
		"username must be at least 3 characters long",
	)
	ErrUsernameTooLong = errors.New(
		"username must not exceed 31 characters",
	)
	ErrUsernameInvalidCharacter = errors.New(
		"username must contain only letters, numbers, and underscores",
	)
)

type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	length := utf8.RuneCountInString(value)

	if length < 3 {
		return Username{}, ErrUsernameTooShort
	}
	if length > 31 {
		return Username{}, ErrUsernameTooLong
	}

	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return Username{}, ErrUsernameInvalidCharacter
		}
	}

	return Username{
		value: value,
	}, nil
}

func (u Username) Value() string {
	return u.value
}

func (u Username) String() string {
	return u.value
}
