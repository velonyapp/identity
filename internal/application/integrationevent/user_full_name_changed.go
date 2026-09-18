package integrationevent

import "time"

type UserFullNameChanged struct {
	BaseIntegrationEvent

	OldFullName string
	NewFullName string
	UpdateTime  time.Time
}

func NewUserFullNameChanged(
	userID string,
	oldFullName string,
	newFullName string,
	updateTime time.Time,
) UserFullNameChanged {
	return UserFullNameChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

		OldFullName: oldFullName,
		NewFullName: newFullName,
		UpdateTime:  updateTime,
	}
}

func (e UserFullNameChanged) Type() string {
	return "user.full-name.changed"
}

func (e UserFullNameChanged) AggregateType() string {
	return "user"
}
