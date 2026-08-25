package vo

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrStorageKeyEmpty = errors.New(
		"storage key must not be empty",
	)
	ErrStorageKeyTooLong = errors.New(
		"storage key must not exceed 1024 bytes",
	)
	ErrStorageKeyInvalidUTF8 = errors.New(
		"storage key must contain valid UTF-8",
	)
)

type StorageKey struct {
	value string
}

func NewStorageKey(value string) (StorageKey, error) {
	if value == "" {
		return StorageKey{}, ErrStorageKeyEmpty
	}

	if len(value) > 1024 {
		return StorageKey{}, ErrStorageKeyTooLong
	}

	if !utf8.ValidString(value) {
		return StorageKey{}, ErrStorageKeyInvalidUTF8
	}

	return StorageKey{
		value: value,
	}, nil
}

func (r StorageKey) Value() string {
	return r.value
}

func (r StorageKey) String() string {
	return r.value
}
