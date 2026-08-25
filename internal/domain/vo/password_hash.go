package vo

import "errors"

var (
	ErrPasswordHashEmpty = errors.New(
		"password hash must not be empty",
	)
)

type PasswordHash struct {
	value string
}

func NewPasswordHash(value string) (PasswordHash, error) {
	if value == "" {
		return PasswordHash{}, ErrPasswordHashEmpty
	}

	return PasswordHash{
		value: value,
	}, nil
}

func (h PasswordHash) Value() string {
	return h.value
}

func (h PasswordHash) String() string {
	return "[REDACTED]"
}
