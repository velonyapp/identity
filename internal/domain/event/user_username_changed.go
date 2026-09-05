package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserUsernameChanged struct {
	BaseDomainEvent

	Username   vo.Username
	UpdateTime vo.Time
}

func NewUserUsernameChanged(
	userID vo.UserID,
	username vo.Username,
	updateTime vo.Time,
) UserUsernameChanged {
	return UserUsernameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		Username:   username,
		UpdateTime: updateTime,
	}
}

func (e UserUsernameChanged) Type() string {
	return "user.username.changed"
}
