package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserUsernameChanged struct {
	BaseDomainEvent

	OldUsername vo.Username
	NewUsername vo.Username
	UpdateTime  vo.Time
}

func NewUserUsernameChanged(
	userID vo.UserID,
	oldUsername vo.Username,
	newUsername vo.Username,
	updateTime vo.Time,
) UserUsernameChanged {
	return UserUsernameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		OldUsername: oldUsername,
		NewUsername: newUsername,
		UpdateTime:  updateTime,
	}
}
