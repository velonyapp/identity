package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserDeleted)(nil)

type UserDeleted struct {
	BaseDomainEvent
}

func NewUserDeleted(
	userID vo.UserID,
	occurTime time.Time,
) *UserDeleted {
	return &UserDeleted{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),
	}
}
