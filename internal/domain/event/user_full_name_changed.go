package event

import "github.com/velony-app/identity/internal/domain/vo"

type UserFullNameChanged struct {
	BaseDomainEvent

	FullName   vo.FullName
	UpdateTime vo.Time
}

func NewUserFullNameChanged(
	userID vo.UserID,
	fullName vo.FullName,
	updateTime vo.Time,
) UserFullNameChanged {
	return UserFullNameChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		FullName:   fullName,
		UpdateTime: updateTime,
	}
}

func (e UserFullNameChanged) Type() string {
	return "user.full-name.changed"
}
