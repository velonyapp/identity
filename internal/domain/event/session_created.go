package event

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ DomainEvent = (*SessionCreated)(nil)

type SessionCreated struct {
	BaseDomainEvent

	userID vo.UserID
}

func NewSessionCreated(
	sessionID vo.SessionID,
	userID vo.UserID,
	occurTime time.Time,
) *SessionCreated {
	return &SessionCreated{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.Value(), occurTime),

		userID: userID,
	}
}

func (e *SessionCreated) UserID() vo.UserID {
	return e.userID
}
