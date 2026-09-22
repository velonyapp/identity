package vo

import "github.com/google/uuid"

type EmailChangeRequestID struct {
	value string
}

func NewEmailChangeRequestID(value string) EmailChangeRequestID {
	return EmailChangeRequestID{value: value}
}

func NewEmailChangeRequestIDRandom() EmailChangeRequestID {
	return EmailChangeRequestID{value: uuid.Must(uuid.NewV7()).String()}
}

func (id EmailChangeRequestID) Value() string {
	return id.value
}

func (id EmailChangeRequestID) String() string {
	return id.value
}
