package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserEmailChanged struct {
	BaseDomainEvent

	OldEmail   *vo.Email
	NewEmail   *vo.Email
	UpdateTime vo.Time
}

func NewUserEmailChanged(
	userID vo.UserID,
	oldEmail *vo.Email,
	newEmail *vo.Email,
	updateTime vo.Time,
) UserEmailChanged {
	return UserEmailChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		OldEmail:   oldEmail,
		NewEmail:   newEmail,
		UpdateTime: updateTime,
	}
}
