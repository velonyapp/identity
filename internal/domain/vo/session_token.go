package vo

import "errors"

var (
	ErrSessionTokenEmpty = errors.New(
		"session token must not be empty",
	)
)

type SessionToken struct {
	value string
}

func NewSessionToken(value string) (SessionToken, error) {
	if value == "" {
		return SessionToken{}, ErrSessionTokenEmpty
	}

	return SessionToken{
		value: value,
	}, nil
}

func (h SessionToken) Value() string {
	return h.value
}

func (h SessionToken) String() string {
	return "[REDACTED]"
}
