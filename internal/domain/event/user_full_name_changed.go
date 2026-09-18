package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserFullNameChanged struct {
	BaseDomainEvent

	OldFullName vo.FullName
	NewFullName vo.FullName
	UpdateTime  vo.Time
}

func NewUserFullNameChanged(
	userID vo.UserID,
	oldFullName vo.FullName,
	newFullName vo.FullName,
	updateTime vo.Time,
) UserFullNameChanged {
	return UserFullNameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		OldFullName: oldFullName,
		NewFullName: newFullName,
		UpdateTime:  updateTime,
	}
}
