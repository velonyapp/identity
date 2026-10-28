package integrationevent

import "time"

var _ IntegrationEvent = (*UserCreated)(nil)

type UserCreated struct {
	BaseIntegrationEvent

	username string
	version  uint32
}

func NewUserCreated(
	aggregateID string,
	occurTime time.Time,
	username string,
	version uint32,
) UserCreated {
	return UserCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		username: username,
		version:  version,
	}
}

func (e UserCreated) Type() string {
	return "identity.user.created"
}

func (e UserCreated) AggregateType() string {
	return "user"
}

func (e UserCreated) Username() string {
	return e.username
}

func (e UserCreated) Version() uint32 {
	return e.version
}
