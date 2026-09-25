package integrationevent

import "time"

var _ IntegrationEvent = (*UserCreated)(nil)

type UserCreated struct {
	BaseIntegrationEvent

	Username string
	FullName string
}

func NewUserCreated(
	userID string,
	username string,
	fullName string,
	occurTime time.Time,
) *UserCreated {
	return &UserCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		Username: username,
		FullName: fullName,
	}
}

func (e UserCreated) Type() string {
	return "user.created"
}

func (e UserCreated) AggregateType() string {
	return "user"
}
