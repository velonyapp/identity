package entity

import (
	"github.com/velonyapp/identity/internal/domain/vo"
)

type Avatar struct {
	id     vo.AvatarID
	userID vo.UserID
	key    vo.AvatarKey
}

func NewAvatar(
	id vo.AvatarID,
	userID vo.UserID,
	key vo.AvatarKey,
) *Avatar {
	return &Avatar{
		id:     id,
		userID: userID,
		key:    key,
	}
}

func ReconstituteAvatar(
	id vo.AvatarID,
	userID vo.UserID,
	key vo.AvatarKey,
) *Avatar {
	return &Avatar{
		id:     id,
		userID: userID,
		key:    key,
	}
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
