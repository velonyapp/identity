package integrationevent

import "time"

type UserUsernameChanged struct {
	BaseIntegrationEvent

	OldUsername string
	NewUsername string
	UpdateTime  time.Time
}

func NewUserUsernameChanged(
	userID string,
	oldUsername string,
	newUsername string,
	updateTime time.Time,
) UserUsernameChanged {
	return UserUsernameChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

		OldUsername: oldUsername,
		NewUsername: newUsername,
		UpdateTime:  updateTime,
	}
}

func (e UserUsernameChanged) Type() string {
	return "user.username.changed"
}

func (e UserUsernameChanged) AggregateType() string {
	return "user"
}
