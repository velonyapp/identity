package event

import "github.com/velonyapp/identity/internal/domain/vo"

type SessionRevoked struct {
	BaseDomainEvent

	RevokeTime vo.Time
}

func NewSessionRevoked(
	sessionID vo.SessionID,
	revokeTime vo.Time,
) SessionRevoked {
	return SessionRevoked{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.String()),

		RevokeTime: revokeTime,
	}
}
