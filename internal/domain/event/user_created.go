package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserCreated)(nil)

type UserCreated struct {
	BaseDomainEvent

	username   vo.Username
	createTime time.Time
	version    uint32
}

func NewUserCreated(
	userID vo.UserID,
	username vo.Username,
	createTime time.Time,
	version uint32,
) UserCreated {
	return UserCreated{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		username:   username,
		createTime: createTime,
		version:    version,
	}
}

func (e UserCreated) Username() vo.Username {
	return e.username
}

func (e UserCreated) CreateTime() time.Time {
	return e.createTime
}

func (e UserCreated) Version() uint32 {
	return e.version
}
