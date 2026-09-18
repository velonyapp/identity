package integrationevent

import "time"

type UserDeleted struct {
	BaseIntegrationEvent

	DeleteTime time.Time
}

func NewUserDeleted(
	userID string,
	deleteTime time.Time,
) UserDeleted {
	return UserDeleted{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

		DeleteTime: deleteTime,
	}
}

func (e UserDeleted) Type() string {
	return "user.deleted"
}

func (e UserDeleted) AggregateType() string {
	return "user"
}
