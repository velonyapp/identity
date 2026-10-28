package integrationevent

import "time"

var _ IntegrationEvent = (*UserDeleted)(nil)

type UserDeleted struct {
	BaseIntegrationEvent

	version uint32
}

func NewUserDeleted(
	aggregateID string,
	occurTime time.Time,
	version uint32,
) UserDeleted {
	return UserDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		version: version,
	}
}

func (e UserDeleted) Type() string {
	return "identity.user.deleted"
}

func (e UserDeleted) AggregateType() string {
	return "user"
}

func (e UserDeleted) Version() uint32 {
	return e.version
}
