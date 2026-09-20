package integrationevent

import "time"

type SessionRevoked struct {
	BaseIntegrationEvent

	RevokeTime time.Time
}

func NewSessionRevoked(
	sessionID string,
	revokeTime time.Time,
) SessionRevoked {
	return SessionRevoked{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID),

		RevokeTime: revokeTime,
	}
}

func (e SessionRevoked) Type() string {
	return "session.revoked"
}

func (e SessionRevoked) AggregateType() string {
	return "session"
}
