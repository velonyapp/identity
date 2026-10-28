package vo

import (
	"errors"

	"github.com/google/uuid"
)

type SessionID struct {
	value string
}

var (
	ErrSessionIDInvalid = errors.New("invalid session ID")
)

func NewSessionID(value string) (SessionID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return SessionID{}, ErrSessionIDInvalid
	}

	return SessionID{value: id.String()}, nil
}

func GenerateSessionID() SessionID {
	return SessionID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id SessionID) String() string {
	return id.value
}

func (id SessionID) Equal(other SessionID) bool {
	return id.value == other.value
}
