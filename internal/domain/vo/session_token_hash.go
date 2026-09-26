package vo

import "errors"

var (
	ErrSessionTokenHashEmpty = errors.New(
		"session token hash must not be empty",
	)
)

type SessionTokenHash struct {
	value []byte
}

func NewSessionTokenHash(value []byte) (SessionTokenHash, error) {
	if len(value) == 0 {
		return SessionTokenHash{}, ErrSessionTokenHashEmpty
	}

	copied := make([]byte, len(value))
	copy(copied, value)

	return SessionTokenHash{
		value: copied,
	}, nil
}

func (h SessionTokenHash) Value() []byte {
	copied := make([]byte, len(h.value))
	copy(copied, h.value)

	return copied
}
