package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserAvatarKeyChanged struct {
	BaseDomainEvent

	OldAvatarKey *vo.AvatarKey
	NewAvatarKey *vo.AvatarKey
}

func NewUserAvatarKeyChanged(
	userID vo.UserID,
	oldAvatarKey *vo.AvatarKey,
	newAvatarKey *vo.AvatarKey,
	occurTime vo.Time,
) UserAvatarKeyChanged {
	return UserAvatarKeyChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),

		OldAvatarKey: oldAvatarKey,
		NewAvatarKey: newAvatarKey,
	}
}

func (e UserAvatarKeyChanged) Type() string {
	return "user.avatar-key.changed"
}

func (e UserAvatarKeyChanged) AggregateType() string {
	return "user"
}
