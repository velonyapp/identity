package vo

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrAvatarKeyEmpty = errors.New(
		"avatar key must not be empty",
	)
	ErrAvatarKeyTooLong = errors.New(
		"avatar key must not exceed 128 bytes",
	)
	ErrAvatarKeyInvalidUTF8 = errors.New(
		"avatar key must contain valid UTF-8",
	)
	ErrAvatarKeyInvalidFormat = errors.New(
		"avatar key must match users/{user_id}/avatar-{token}.webp",
	)
)

type AvatarKey struct {
	userID string
	token  string
}

func NewAvatarKey(value string) (AvatarKey, error) {
	if value == "" {
		return AvatarKey{}, ErrAvatarKeyEmpty
	}

	if len(value) > 128 {
		return AvatarKey{}, ErrAvatarKeyTooLong
	}

	if !utf8.ValidString(value) {
		return AvatarKey{}, ErrAvatarKeyInvalidUTF8
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
		userID: userID,
		token:  token,
	}, nil
}

func (k AvatarKey) UserID() string {
	return k.userID
}

func (k AvatarKey) Token() string {
	return k.token
}

func (k AvatarKey) String() string {
	return "users/" + k.userID + "/avatar-" + k.token + ".webp"
}
