package integrationevent

import "time"

type SessionRefreshed struct {
	BaseIntegrationEvent

	RefreshTime time.Time
}

func NewSessionRefreshed(
	sessionID string,
	refreshTime time.Time,
) SessionRefreshed {
	return SessionRefreshed{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID),

		RefreshTime: refreshTime,
	}
}

func (e SessionRefreshed) Type() string {
	return "session.refreshed"
}

func (e SessionRefreshed) AggregateType() string {
	return "session"
}
