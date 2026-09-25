package integrationevent

import "time"

var _ IntegrationEvent = (*SessionCreated)(nil)

type SessionCreated struct {
	BaseIntegrationEvent

	userID string
}

func NewSessionCreated(
	sessionID string,
	userID string,
	occurTime time.Time,
) *SessionCreated {
	return &SessionCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID, occurTime),

		userID: userID,
	}
}

func (e *SessionCreated) Type() string {
	return "session.created"
}

func (e *SessionCreated) AggregateType() string {
	return "session"
}

func (e *SessionCreated) UserID() string {
	return e.userID
}
