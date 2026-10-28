package vo

import (
	"errors"

	"github.com/google/uuid"
)

type AvatarID struct {
	value string
}

var (
	ErrAvatarIDInvalid = errors.New("invalid avatar ID")
)

func NewAvatarID(value string) (AvatarID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return AvatarID{}, ErrAvatarIDInvalid
	}

	return AvatarID{value: id.String()}, nil
}

func GenerateAvatarID() AvatarID {
	return AvatarID{
		value: uuid.Must(uuid.NewV7()).String(),
	}
}

func (id AvatarID) String() string {
	return id.value
}

func (id AvatarID) Equal(other AvatarID) bool {
	return id.value == other.value
}
