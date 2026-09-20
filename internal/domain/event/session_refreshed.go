package event

import "github.com/velonyapp/identity/internal/domain/vo"

type SessionRefreshed struct {
	BaseDomainEvent

	RefreshTime vo.Time
}

func NewSessionRefreshed(
	sessionID vo.SessionID,
	refreshTime vo.Time,
) SessionRefreshed {
	return SessionRefreshed{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.String()),

		RefreshTime: refreshTime,
	}
}
