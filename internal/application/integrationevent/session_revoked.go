package integrationevent

import "time"

var _ IntegrationEvent = (*SessionRevoked)(nil)

type SessionRevoked struct {
	BaseIntegrationEvent
}

func NewSessionRevoked(
	sessionID string,
	occurTime time.Time,
) *SessionRevoked {
	return &SessionRevoked{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID, occurTime),
	}
}

func (e *SessionRevoked) Type() string {
	return "session.revoked"
}

func (e *SessionRevoked) AggregateType() string {
	return "session"
}
