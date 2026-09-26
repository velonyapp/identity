package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserEmailChangeRequested)(nil)

type UserEmailChangeRequested struct {
	BaseDomainEvent

	email      vo.Email
	expireTime time.Time
}

func NewUserEmailChangeRequested(
	userID vo.UserID,
	email vo.Email,
	expireTime time.Time,
	occurTime time.Time,
) *UserEmailChangeRequested {
	return &UserEmailChangeRequested{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		email:      email,
		expireTime: expireTime,
	}
}

func (e *UserEmailChangeRequested) Email() vo.Email {
	return e.email
}

func (e *UserEmailChangeRequested) ExpireTime() time.Time {
	return e.expireTime
}
