package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserAvatarKeyChanged struct {
	BaseDomainEvent

	OldAvatarKey *vo.AvatarKey
	NewAvatarKey *vo.AvatarKey
	UpdateTime   vo.Time
}

func NewUserAvatarKeyChanged(
	userID vo.UserID,
	oldAvatarKey *vo.AvatarKey,
	newAvatarKey *vo.AvatarKey,
	updateTime vo.Time,
) UserAvatarKeyChanged {
	return UserAvatarKeyChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		OldAvatarKey: oldAvatarKey,
		NewAvatarKey: newAvatarKey,
		UpdateTime:   updateTime,
	}
}
