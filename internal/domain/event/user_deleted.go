package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserDeleted struct {
	BaseDomainEvent
}

func NewUserDeleted(
	userID vo.UserID,
	occurTime vo.Time,
) UserDeleted {
	return UserDeleted{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),
	}
}

func (e UserDeleted) Type() string {
	return "user.deleted"
}

func (e UserDeleted) AggregateType() string {
	return "user"
}
