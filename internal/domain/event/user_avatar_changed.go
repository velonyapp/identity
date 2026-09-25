package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*UserAvatarChanged)(nil)

type UserAvatarChanged struct {
	BaseDomainEvent

	oldAvatarKey *vo.AvatarKey
	newAvatarKey *vo.AvatarKey
}

func NewUserAvatarChanged(
	userID vo.UserID,
	oldAvatarKey *vo.AvatarKey,
	newAvatarKey *vo.AvatarKey,
	occurTime time.Time,
) *UserAvatarChanged {
	return &UserAvatarChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.Value(), occurTime),

		oldAvatarKey: oldAvatarKey,
		newAvatarKey: newAvatarKey,
	}
}

func (e *UserAvatarChanged) OldAvatarKey() *vo.AvatarKey {
	if e.oldAvatarKey == nil {
		return nil
	}
	value := *e.oldAvatarKey
	return &value
}

func (e *UserAvatarChanged) NewAvatarKey() *vo.AvatarKey {
	if e.newAvatarKey == nil {
		return nil
	}
	value := *e.newAvatarKey
	return &value
}
