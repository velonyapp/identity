package vo

import "errors"

var (
	ErrSessionTokenHashEmpty = errors.New(
		"session token hash must not be empty",
	)
)

type SessionTokenHash struct {
	value string
}

func NewSessionTokenHash(value string) (SessionTokenHash, error) {
	if value == "" {
		return SessionTokenHash{}, ErrSessionTokenHashEmpty
	}

	return SessionTokenHash{
		value: value,
	}, nil
}

func (h SessionTokenHash) Value() string {
	return h.value
}
