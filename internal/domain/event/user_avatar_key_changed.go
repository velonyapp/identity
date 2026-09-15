package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserAvatarKeyChanged struct {
	BaseDomainEvent

	AvatarKey  *vo.AvatarKey
	UpdateTime vo.Time
}

func NewUserAvatarKeyChanged(
	userID vo.UserID,
	avatarKey *vo.AvatarKey,
	updateTime vo.Time,
) UserAvatarKeyChanged {
	return UserAvatarKeyChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		AvatarKey:  avatarKey,
		UpdateTime: updateTime,
	}
}

func (e UserAvatarKeyChanged) Type() string {
	return "user.avatar-key.changed"
}
