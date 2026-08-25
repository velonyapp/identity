package vo

import "errors"

var (
	ErrGoogleSubEmpty = errors.New(
		"google sub must not be empty",
	)
	ErrGoogleSubTooLong = errors.New(
		"google sub must not exceed 255 bytes",
	)
	ErrGoogleSubNonASCII = errors.New(
		"google sub must contain only ASCII characters",
	)
)

type GoogleSub struct {
	value string
}

func NewGoogleSub(value string) (GoogleSub, error) {
	if value == "" {
		return GoogleSub{}, ErrGoogleSubEmpty
	}

	if len(value) > 255 {
		return GoogleSub{}, ErrGoogleSubTooLong
	}
	for i := 0; i < len(value); i++ {
		if value[i] > 127 {
			return GoogleSub{}, ErrGoogleSubNonASCII
		}
	}

	return GoogleSub{
		value: value,
	}, nil
}

func (s GoogleSub) Value() string {
	return s.value
}

func (s GoogleSub) String() string {
	return s.value
}
