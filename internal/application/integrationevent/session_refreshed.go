package integrationevent

import "time"

var _ IntegrationEvent = (*SessionRefreshed)(nil)

type SessionRefreshed struct {
	BaseIntegrationEvent
}

func NewSessionRefreshed(
	sessionID string,
	occurTime time.Time,
) *SessionRefreshed {
	return &SessionRefreshed{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID, occurTime),
	}
}

func (e *SessionRefreshed) Type() string {
	return "session.refreshed"
}

func (e *SessionRefreshed) AggregateType() string {
	return "session"
}
