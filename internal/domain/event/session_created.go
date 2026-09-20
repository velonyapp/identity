package event

import "github.com/velonyapp/identity/internal/domain/vo"

type SessionCreated struct {
	BaseDomainEvent

	UserID     vo.UserID
	CreateTime vo.Time
}

func NewSessionCreated(
	sessionID vo.SessionID,
	userID vo.UserID,
	createTime vo.Time,
) SessionCreated {
	return SessionCreated{
		BaseDomainEvent: NewBaseDomainEvent(sessionID.String()),

		UserID:     userID,
		CreateTime: createTime,
	}
}
