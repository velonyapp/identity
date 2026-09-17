package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserUsernameChanged struct {
	BaseDomainEvent

	OldUsername vo.Username
	NewUsername vo.Username
}

func NewUserUsernameChanged(
	userID vo.UserID,
	oldUsername vo.Username,
	newUsername vo.Username,
	occurTime vo.Time,
) UserUsernameChanged {
	return UserUsernameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),

		OldUsername: oldUsername,
		NewUsername: newUsername,
	}
}

func (e UserUsernameChanged) Type() string {
	return "user.username.changed"
}

func (e UserUsernameChanged) AggregateType() string {
	return "user"
}
