package vo

import (
	"bytes"
	"errors"
)

var (
	ErrPasswordHashEmpty = errors.New("password hash must not be empty")
)

type PasswordHash struct {
	value []byte
}

func NewPasswordHash(value []byte) (PasswordHash, error) {
	if len(value) == 0 {
		return PasswordHash{}, ErrPasswordHashEmpty
	}

	return PasswordHash{
		value: bytes.Clone(value),
	}, nil
}

func (h PasswordHash) Bytes() []byte {
	return bytes.Clone(h.value)
}

func (h PasswordHash) Equal(other PasswordHash) bool {
	return bytes.Equal(h.value, other.value)
}
