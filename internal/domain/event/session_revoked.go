package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*SessionRevoked)(nil)

type SessionRevoked struct {
	BaseDomainEvent
}

func NewSessionRevoked(
	sessionID vo.SessionID,
	occurTime time.Time,
) *SessionRevoked {
	return &SessionRevoked{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.Value(), occurTime),
	}
}
