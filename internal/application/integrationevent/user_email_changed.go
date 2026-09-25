package integrationevent

import "time"

var _ IntegrationEvent = (*UserEmailChanged)(nil)

type UserEmailChanged struct {
	BaseIntegrationEvent

	oldEmail *string
	newEmail *string
}

func NewUserEmailChanged(
	userID string,
	oldEmail *string,
	newEmail *string,
	occurTime time.Time,
) *UserEmailChanged {
	return &UserEmailChanged{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		oldEmail: oldEmail,
		newEmail: newEmail,
	}
}

func (e *UserEmailChanged) Type() string {
	return "user.email.changed"
}

func (e *UserEmailChanged) AggregateType() string {
	return "user"
}

func (e *UserEmailChanged) OldEmail() *string {
	if e.oldEmail == nil {
		return nil
	}
	value := *e.oldEmail
	return &value
}

func (e *UserEmailChanged) NewEmail() *string {
	if e.newEmail == nil {
		return nil
	}
	value := *e.newEmail
	return &value
}
