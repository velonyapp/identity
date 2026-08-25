package event

import "github.com/velony-app/identity/internal/domain/vo"

type UserAvatarKeyChanged struct {
	BaseDomainEvent

	AvatarKey  *vo.StorageKey
	UpdateTime vo.Time
}

func NewUserAvatarKeyChanged(
	userID vo.UserID,
	avatarKey *vo.StorageKey,
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
