package integrationevent

import "time"

var _ IntegrationEvent = (*UserAvatarChanged)(nil)

type UserAvatarChanged struct {
	BaseIntegrationEvent

	oldAvatarKey *string
	newAvatarKey *string
}

func NewUserAvatarChanged(
	userID string,
	oldAvatarKey *string,
	newAvatarKey *string,
	occurTime time.Time,
) *UserAvatarChanged {
	return &UserAvatarChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		oldAvatarKey: oldAvatarKey,
		newAvatarKey: newAvatarKey,
	}
}

func (e *UserAvatarChanged) Type() string {
	return "user.avatar-key.changed"
}

func (e *UserAvatarChanged) AggregateType() string {
	return "user"
}

func (e *UserAvatarChanged) OldAvatarKey() *string {
	if e.oldAvatarKey == nil {
		return nil
	}
	value := *e.oldAvatarKey
	return &value
}

func (e *UserAvatarChanged) NewAvatarKey() *string {
	if e.newAvatarKey == nil {
		return nil
	}
	value := *e.newAvatarKey
	return &value
}
