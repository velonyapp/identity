package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserDeleted)(nil)

type UserDeleted struct {
	BaseDomainEvent

	deleteTime time.Time
	version    uint32
}

func NewUserDeleted(
	userID vo.UserID,
	deleteTime time.Time,
	version uint32,
) UserDeleted {
	return UserDeleted{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),
		deleteTime:      deleteTime,
		version:         version,
	}
}

func (e UserDeleted) DeleteTime() time.Time {
	return e.deleteTime
}

func (e UserDeleted) Version() uint32 {
	return e.version
}
