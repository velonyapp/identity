package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserEmailChanged struct {
	BaseDomainEvent

	Email      *vo.Email
	UpdateTime vo.Time
}

func NewUserEmailChanged(
	userID vo.UserID,
	avatarKey *vo.Email,
	updateTime vo.Time,
) UserEmailChanged {
	return UserEmailChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		Email:      avatarKey,
		UpdateTime: updateTime,
	}
}

func (e UserEmailChanged) Type() string {
	return "user.email.changed"
}
