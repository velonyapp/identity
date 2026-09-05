package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserDeleted struct {
	BaseDomainEvent

	DeleteTime vo.Time
}

func NewUserDeleted(
	userID vo.UserID,
	deleteTime vo.Time,
) UserDeleted {
	return UserDeleted{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		DeleteTime: deleteTime,
	}
}

func (e UserDeleted) Type() string {
	return "user.deleted"
}
