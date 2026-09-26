package integrationevent

import "time"

var _ IntegrationEvent = (*UserEmailChangeRequested)(nil)

type UserEmailChangeRequested struct {
	BaseIntegrationEvent

	email      string
	expireTime time.Time
}

func NewUserEmailChangeRequested(
	userID string,
	email string,
	expireTime time.Time,
	occurTime time.Time,
) *UserEmailChangeRequested {
	return &UserEmailChangeRequested{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID, occurTime),

		email:      email,
		expireTime: expireTime,
	}
}

func (e *UserEmailChangeRequested) Type() string {
	return "user.email.change-requested"
}

func (e *UserEmailChangeRequested) AggregateType() string {
	return "user"
}

func (e *UserEmailChangeRequested) Email() string {
	return e.email
}

func (e *UserEmailChangeRequested) ExpireTime() time.Time {
	return e.expireTime
}
