package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserEmailChanged)(nil)

type UserEmailChanged struct {
	BaseDomainEvent

	oldEmail *vo.Email
	newEmail *vo.Email
}

func NewUserEmailChanged(
	userID vo.UserID,
	oldEmail *vo.Email,
	newEmail *vo.Email,
	occurTime time.Time,
) *UserEmailChanged {
	return &UserEmailChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		oldEmail: oldEmail,
		newEmail: newEmail,
	}
}

func (e *UserEmailChanged) OldEmail() *vo.Email {
	if e.oldEmail == nil {
		return nil
	}
	value := *e.oldEmail
	return &value
}

func (e *UserEmailChanged) NewEmail() *vo.Email {
	if e.newEmail == nil {
		return nil
	}
	value := *e.newEmail
	return &value
}
