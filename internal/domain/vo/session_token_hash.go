package vo

import (
	"bytes"
	"crypto/subtle"
	"errors"
)

var (
	ErrSessionTokenHashEmpty = errors.New("session token hash must not be empty")
)

type SessionTokenHash struct {
	value []byte
}

func NewSessionTokenHash(value []byte) (SessionTokenHash, error) {
	if len(value) == 0 {
		return SessionTokenHash{}, ErrSessionTokenHashEmpty
	}

	return SessionTokenHash{
		value: bytes.Clone(value),
	}, nil
}

func (h SessionTokenHash) Bytes() []byte {
	return bytes.Clone(h.value)
}

func (h SessionTokenHash) Equal(other SessionTokenHash) bool {
	return subtle.ConstantTimeCompare(h.value, other.value) == 1
}
