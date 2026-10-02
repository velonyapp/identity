package integrationevent

import "time"

var _ IntegrationEvent = (*UserAvatarChanged)(nil)

type UserAvatarChanged struct {
	BaseIntegrationEvent

	oldAvatarID *string
	newAvatarID *string
}

func NewUserAvatarChanged(
	userID string,
	oldAvatarID *string,
	newAvatarID *string,
	occurTime time.Time,
) *UserAvatarChanged {
	return &UserAvatarChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		oldAvatarID: oldAvatarID,
		newAvatarID: newAvatarID,
	}
}

func (e *UserAvatarChanged) Type() string {
	return "identity.user.avatar.changed"
}

func (e *UserAvatarChanged) AggregateType() string {
	return "user"
}

func (e *UserAvatarChanged) OldAvatarID() *string {
	if e.oldAvatarID == nil {
		return nil
	}
	value := *e.oldAvatarID
	return &value
}

func (e *UserAvatarChanged) NewAvatarID() *string {
	if e.newAvatarID == nil {
		return nil
	}
	value := *e.newAvatarID
	return &value
}
