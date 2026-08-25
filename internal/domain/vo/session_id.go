package vo

import "github.com/google/uuid"

type SessionID struct {
	value string
}

func NewSessionID(value string) SessionID {
	return SessionID{value: value}
}

func NewSessionIDRandom() SessionID {
	return SessionID{value: uuid.Must(uuid.NewV7()).String()}
}

func (id SessionID) Value() string {
	return id.value
}

func (id SessionID) String() string {
	return id.value
}
