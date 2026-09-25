package integrationevent

import "time"

var _ IntegrationEvent = (*UserFullNameChanged)(nil)

type UserFullNameChanged struct {
	BaseIntegrationEvent

	oldFullName string
	newFullName string
}

func NewUserFullNameChanged(
	userID string,
	oldFullName string,
	newFullName string,
	occurTime time.Time,
) *UserFullNameChanged {
	return &UserFullNameChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		oldFullName: oldFullName,
		newFullName: newFullName,
	}
}

func (e *UserFullNameChanged) Type() string {
	return "user.full-name.changed"
}

func (e *UserFullNameChanged) AggregateType() string {
	return "user"
}

func (e *UserFullNameChanged) OldFullName() string {
	return e.oldFullName
}

func (e *UserFullNameChanged) NewFullName() string {
	return e.newFullName
}
