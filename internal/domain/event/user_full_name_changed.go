package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserFullNameChanged)(nil)

type UserFullNameChanged struct {
	BaseDomainEvent

	oldFullName vo.FullName
	newFullName vo.FullName
}

func NewUserFullNameChanged(
	userID vo.UserID,
	oldFullName vo.FullName,
	newFullName vo.FullName,
	occurTime time.Time,
) *UserFullNameChanged {
	return &UserFullNameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		oldFullName: oldFullName,
		newFullName: newFullName,
	}
}

func (e *UserFullNameChanged) OldFullName() vo.FullName {
	return e.oldFullName
}

func (e *UserFullNameChanged) NewFullName() vo.FullName {
	return e.newFullName
}
