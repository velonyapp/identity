package vo

import "errors"

var (
	ErrRequestTokenEmpty = errors.New(
		"request token must not be empty",
	)
)

type RequestToken struct {
	value string
}

func NewRequestToken(value string) (RequestToken, error) {
	if value == "" {
		return RequestToken{}, ErrRequestTokenEmpty
	}

	return RequestToken{
		value: value,
	}, nil
}

func (h RequestToken) Value() string {
	return h.value
}

func (h RequestToken) String() string {
	return "[REDACTED]"
}
