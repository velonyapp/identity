package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserFullNameChanged struct {
	BaseDomainEvent

	OldFullName vo.FullName
	NewFullName vo.FullName
}

func NewUserFullNameChanged(
	userID vo.UserID,
	oldFullName vo.FullName,
	newFullName vo.FullName,
	occurTime vo.Time,
) UserFullNameChanged {
	return UserFullNameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),

		OldFullName: oldFullName,
		NewFullName: newFullName,
	}
}

func (e UserFullNameChanged) Type() string {
	return "user.full-name.changed"
}

func (e UserFullNameChanged) AggregateType() string {
	return "user"
}
