package integrationevent

import "time"

var _ IntegrationEvent = (*UserUsernameChanged)(nil)

type UserUsernameChanged struct {
	BaseIntegrationEvent

	oldUsername string
	newUsername string
}

func NewUserUsernameChanged(
	userID string,
	oldUsername string,
	newUsername string,
	occurTime time.Time,
) *UserUsernameChanged {
	return &UserUsernameChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		oldUsername: oldUsername,
		newUsername: newUsername,
	}
}

func (e *UserUsernameChanged) Type() string {
	return "user.username.changed"
}

func (e *UserUsernameChanged) AggregateType() string {
	return "user"
}

func (e *UserUsernameChanged) OldUsername() string {
	return e.oldUsername
}

func (e *UserUsernameChanged) NewUsername() string {
	return e.newUsername
}
