package vo

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrAvatarKeyEmpty            = errors.New("avatar key must not be empty")
	ErrAvatarKeyTooLong          = errors.New("avatar key must not exceed 128 bytes")
	ErrAvatarKeyInvalidCharacter = errors.New("avatar key contains an invalid character")
	ErrAvatarKeyInvalidFormat    = errors.New("avatar key must match users/{user_id}/avatar-{token}.webp")
)

type AvatarKey struct {
	value string
}

func NewAvatarKey(value string) (AvatarKey, error) {
	if value == "" {
		return AvatarKey{}, ErrAvatarKeyEmpty
	}

	for _, c := range value {
		if !(unicode.IsLetter(c) ||
			unicode.IsNumber(c) ||
			c == '-' || c == '_' || c == '.' || c == '/') {
			return AvatarKey{}, ErrUsernameInvalidCharacter
		}
	}

	length := utf8.RuneCountInString(value)

	if length > 128 {
		return AvatarKey{}, ErrAvatarKeyTooLong
	}

	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return AvatarKey{}, ErrAvatarKeyInvalidFormat
	}

	if parts[0] != "users" {
		return AvatarKey{}, ErrAvatarKeyInvalidFormat
	}

	userID := parts[1]
	if userID == "" {
		return AvatarKey{}, ErrAvatarKeyInvalidFormat
	}

	fileName := parts[2]

	if !strings.HasPrefix(fileName, "avatar-") ||
		!strings.HasSuffix(fileName, ".webp") {
		return AvatarKey{}, ErrAvatarKeyInvalidFormat
	}

	token := fileName[len("avatar-") : len(fileName)-len(".webp")]
	if token == "" {
		return AvatarKey{}, ErrAvatarKeyInvalidFormat
	}

	return AvatarKey{
		value: value,
	}, nil
}

func (k AvatarKey) String() string {
	return k.value
}

func (k AvatarKey) Equal(other AvatarKey) bool {
	return k.value == other.value
}
