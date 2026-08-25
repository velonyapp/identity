package vo

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrFullNameTooShort = errors.New(
		"full name must be at least 3 characters long",
	)
	ErrFullNameTooLong = errors.New(
		"full name must not exceed 63 characters",
	)
)

type FullName struct {
	value string
}

func NewFullName(value string) (FullName, error) {
	value = strings.TrimSpace(value)

	length := utf8.RuneCountInString(value)

	if length < 3 {
		return FullName{}, ErrFullNameTooShort
	}

	if length > 63 {
		return FullName{}, ErrFullNameTooLong
	}

	return FullName{
		value: value,
	}, nil
}

func (n FullName) Value() string {
	return n.value
}

func (n FullName) String() string {
	return n.value
}
