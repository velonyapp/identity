package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserAvatarChanged)(nil)

type UserAvatarChanged struct {
	BaseDomainEvent

	oldAvatarID *vo.AvatarID
	newAvatarID *vo.AvatarID
}

func NewUserAvatarChanged(
	userID vo.UserID,
	oldAvatarID *vo.AvatarID,
	newAvatarID *vo.AvatarID,
	occurTime time.Time,
) *UserAvatarChanged {
	return &UserAvatarChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		oldAvatarID: oldAvatarID,
		newAvatarID: newAvatarID,
	}
}

func (e *UserAvatarChanged) OldAvatarID() *vo.AvatarID {
	if e.oldAvatarID == nil {
		return nil
	}
	value := *e.oldAvatarID
	return &value
}

func (e *UserAvatarChanged) NewAvatarID() *vo.AvatarID {
	if e.newAvatarID == nil {
		return nil
	}
	value := *e.newAvatarID
	return &value
}
