package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserUpdated)(nil)

type UserUpdated struct {
	BaseDomainEvent

	username   vo.Username
	email      *vo.Email
	updateTime time.Time
	version    uint32
}

func NewUserUpdated(
	userID vo.UserID,
	username vo.Username,
	email *vo.Email,
	updateTime time.Time,
	version uint32,
) UserUpdated {
	u := UserUpdated{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		username:   username,
		updateTime: updateTime,
		version:    version,
	}

	if email != nil {
		value := *email
		u.email = &value
	}

	return u
}

func (e UserUpdated) Username() vo.Username {
	return e.username
}

func (e UserUpdated) Email() *vo.Email {
	if e.email == nil {
		return nil
	}

	value := *e.email
	return &value
}

func (e UserUpdated) UpdateTime() time.Time {
	return e.updateTime
}

func (e UserUpdated) Version() uint32 {
	return e.version
}
