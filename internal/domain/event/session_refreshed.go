package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*SessionRefreshed)(nil)

type SessionRefreshed struct {
	BaseDomainEvent
}

func NewSessionRefreshed(
	sessionID vo.SessionID,
	occurTime time.Time,
) *SessionRefreshed {
	return &SessionRefreshed{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.Value(), occurTime),
	}
}
