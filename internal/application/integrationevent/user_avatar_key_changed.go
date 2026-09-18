package integrationevent

import "time"

type UserAvatarKeyChanged struct {
	BaseIntegrationEvent

	OldAvatarKey *string
	NewAvatarKey *string
	UpdateTime   time.Time
}

func NewUserAvatarKeyChanged(
	userID string,
	oldAvatarKey *string,
	newAvatarKey *string,
	updateTime time.Time,
) UserAvatarKeyChanged {
	return UserAvatarKeyChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

		OldAvatarKey: oldAvatarKey,
		NewAvatarKey: newAvatarKey,
		UpdateTime:   updateTime,
	}
}

func (e UserAvatarKeyChanged) Type() string {
	return "user.avatar-key.changed"
}

func (e UserAvatarKeyChanged) AggregateType() string {
	return "user"
}
