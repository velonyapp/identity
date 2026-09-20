package integrationevent

import "time"

type SessionCreated struct {
	BaseIntegrationEvent

	UserID     string
	CreateTime time.Time
}

func NewSessionCreated(
	sessionID string,
	userID string,
	createTime time.Time,
) SessionCreated {
	return SessionCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(sessionID),

		UserID:     userID,
		CreateTime: createTime,
	}
}

func (e SessionCreated) Type() string {
	return "session.created"
}

func (e SessionCreated) AggregateType() string {
	return "session"
}
