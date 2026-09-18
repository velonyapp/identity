package integrationevent

import "time"

type UserEmailChanged struct {
	BaseIntegrationEvent

	OldEmail   *string
	NewEmail   *string
	UpdateTime time.Time
}

func NewUserEmailChanged(
	userID string,
	oldEmail *string,
	newEmail *string,
	updateTime time.Time,
) UserEmailChanged {
	return UserEmailChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

		OldEmail:   oldEmail,
		NewEmail:   newEmail,
		UpdateTime: updateTime,
	}
}

func (e UserEmailChanged) Type() string {
	return "user.email.changed"
}

func (e UserEmailChanged) AggregateType() string {
	return "user"
}
