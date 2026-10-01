package integrationevent

import "time"

var _ IntegrationEvent = (*UserDeleted)(nil)

type UserDeleted struct {
	BaseIntegrationEvent
}

func NewUserDeleted(
	userID string,
	occurTime time.Time,
) *UserDeleted {
	return &UserDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),
	}
}

func (e *UserDeleted) Type() string {
	return "identity.user.deleted"
}

func (e *UserDeleted) AggregateType() string {
	return "user"
}
