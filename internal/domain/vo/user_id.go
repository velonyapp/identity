package vo

import (
	"errors"

	"github.com/google/uuid"
)

type UserID struct {
	value string
}

var (
	ErrUserIDInvalid = errors.New("invalid user ID")
)

func NewUserID(value string) (UserID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, ErrUserIDInvalid
	}

	return UserID{value: id.String()}, nil
}

func GenerateUserID() UserID {
	return UserID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id UserID) String() string {
	return id.value
}

func (id UserID) Equal(other UserID) bool {
	return id.value == other.value
}
