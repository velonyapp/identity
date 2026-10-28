package vo

import (
	"errors"
	"unicode"
)

var (
	ErrEmailChangeRequestTokenInvalid = errors.New("email change request token must be exactly 6 digits")
)

type EmailChangeRequestToken struct {
	value string
}

func NewEmailChangeRequestToken(value string) (EmailChangeRequestToken, error) {
	if len(value) != 6 {
		return EmailChangeRequestToken{}, ErrEmailChangeRequestTokenInvalid
	}

	for _, c := range value {
		if !unicode.IsDigit(c) {
			return EmailChangeRequestToken{}, ErrEmailChangeRequestTokenInvalid
		}
	}

	return EmailChangeRequestToken{
		value: value,
	}, nil
}

func (h EmailChangeRequestToken) String() string {
	return h.value
}

func (h EmailChangeRequestToken) Bytes() []byte {
	return []byte(h.value)
}

func (h EmailChangeRequestToken) Equal(other EmailChangeRequestToken) bool {
	return h.value == other.value
}
