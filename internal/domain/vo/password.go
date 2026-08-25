package vo

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrPasswordTooShort = errors.New(
		"password must be at least 8 characters long",
	)
	ErrPasswordTooLong = errors.New(
		"password must not exceed 64 characters",
	)
	ErrPasswordMissingUppercase = errors.New(
		"password must contain at least one uppercase character",
	)
	ErrPasswordMissingLowercase = errors.New(
		"password must contain at least one lowercase character",
	)
	ErrPasswordMissingNumber = errors.New(
		"password must contain at least one number",
	)
	ErrPasswordMissingSymbol = errors.New(
		"password must contain at least one symbol",
	)
	ErrPasswordContainsWhitespace = errors.New(
		"password must not contain whitespace",
	)
)

type Password struct {
	value string
}

func NewPassword(value string) (Password, error) {
	value = strings.TrimSpace(value)

	length := utf8.RuneCountInString(value)

	if length < 8 {
		return Password{}, ErrPasswordTooShort
	}
	if length > 64 {
		return Password{}, ErrPasswordTooLong
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
		hasSymbol bool
	)

	for _, char := range value {
		switch {
		case unicode.IsSpace(char):
			return Password{}, ErrPasswordContainsWhitespace

		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasLower = true

		case unicode.IsDigit(char):
			hasNumber = true

		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return Password{}, ErrPasswordMissingUppercase
	}
	if !hasLower {
		return Password{}, ErrPasswordMissingLowercase
	}
	if !hasNumber {
		return Password{}, ErrPasswordMissingNumber
	}
	if !hasSymbol {
		return Password{}, ErrPasswordMissingSymbol
	}

	return Password{
		value: value,
	}, nil
}

func (p Password) Value() string {
	return p.value
}

func (p Password) String() string {
	return "[REDACTED]"
}
