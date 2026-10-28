package integrationevent

import "time"

var _ IntegrationEvent = (*UserUpdated)(nil)

type UserUpdated struct {
	BaseIntegrationEvent

	username string
	email    *string
	version  uint32
}

func NewUserUpdated(
	aggregateID string,
	occurTime time.Time,
	username string,
	email *string,
	version uint32,
) UserUpdated {
	e := UserUpdated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(aggregateID, occurTime),

		username: username,
		version:  version,
	}

	if email != nil {
		value := *email
		e.email = &value
	}

	return e
}

func (e UserUpdated) Type() string {
	return "identity.user.updated"
}

func (e UserUpdated) AggregateType() string {
	return "user"
}

func (e UserUpdated) Username() string {
	return e.username
}

func (e UserUpdated) Email() *string {
	if e.email == nil {
		return nil
	}

	value := *e.email
	return &value
}

func (e UserUpdated) Version() uint32 {
	return e.version
}
