package entity

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

type Avatar struct {
	id         vo.AvatarID
	userID     vo.UserID
	key        vo.AvatarKey
	deleteTime *time.Time
}

func ReconstituteAvatar(
	id vo.AvatarID,
	userID vo.UserID,
	key vo.AvatarKey,
	deleteTime *time.Time,
) *Avatar {
	a := &Avatar{
		id:     id,
		userID: userID,
		key:    key,
	}

	if deleteTime != nil {
		value := *deleteTime
		a.deleteTime = &value
	}

	return a
}

func (a *Avatar) ID() vo.AvatarID {
	return a.id
}

func (a *Avatar) UserID() vo.UserID {
	return a.userID
}

func (a *Avatar) Key() vo.AvatarKey {
	return a.key
}

func (a *Avatar) IsDeleted() bool {
	return a.deleteTime != nil
}

func (a *Avatar) Delete(now time.Time) {
	if a.IsDeleted() {
		return
	}

	a.deleteTime = &now
}
