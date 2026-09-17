package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserEmailChanged struct {
	BaseDomainEvent

	OldEmail *vo.Email
	NewEmail *vo.Email
}

func NewUserEmailChanged(
	userID vo.UserID,
	oldEmail *vo.Email,
	newEmail *vo.Email,
	occurTime vo.Time,
) UserEmailChanged {
	return UserEmailChanged{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),

		OldEmail: oldEmail,
		NewEmail: newEmail,
	}
}

func (e UserEmailChanged) Type() string {
	return "user.email.changed"
}

func (e UserEmailChanged) AggregateType() string {
	return "user"
}
