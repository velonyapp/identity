package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserUsernameChanged)(nil)

type UserUsernameChanged struct {
	BaseDomainEvent

	oldUsername vo.Username
	newUsername vo.Username
}

func NewUserUsernameChanged(
	userID vo.UserID,
	oldUsername vo.Username,
	newUsername vo.Username,
	occurTime time.Time,
) *UserUsernameChanged {
	return &UserUsernameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		oldUsername: oldUsername,
		newUsername: newUsername,
	}
}

func (e *UserUsernameChanged) OldUsername() vo.Username {
	return e.oldUsername
}

func (e *UserUsernameChanged) NewUsername() vo.Username {
	return e.newUsername
}
