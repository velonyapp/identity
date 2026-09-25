package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserCreated)(nil)

type UserCreated struct {
	BaseDomainEvent

	username vo.Username
	fullName vo.FullName
}

func NewUserCreated(
	userID vo.UserID,
	username vo.Username,
	fullName vo.FullName,
	occurTime time.Time,
) *UserCreated {
	return &UserCreated{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		username: username,
		fullName: fullName,
	}
}

func (e *UserCreated) Username() vo.Username {
	return e.username
}

func (e *UserCreated) FullName() vo.FullName {
	return e.fullName
}
