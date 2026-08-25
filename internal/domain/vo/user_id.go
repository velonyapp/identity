package vo

import "github.com/google/uuid"

type UserID struct {
	value string
}

func NewUserID(value string) UserID {
	return UserID{value: value}
}

func NewUserIDRandom() UserID {
	return UserID{value: uuid.Must(uuid.NewV7()).String()}
}

func (id UserID) Value() string {
	return id.value
}

func (id UserID) String() string {
	return id.value
}
