package event

import "github.com/velonyapp/identity/internal/domain/vo"

type UserCreated struct {
	BaseDomainEvent

	Username  vo.Username
	Email     *vo.Email
	FullName  vo.FullName
	AvatarKey *vo.AvatarKey
}

func NewUserCreated(
	userID vo.UserID,
	username vo.Username,
	email *vo.Email,
	fullName vo.FullName,
	avatarKey *vo.AvatarKey,
	occurTime vo.Time,
) UserCreated {
	return UserCreated{
		BaseDomainEvent: NewBaseDomainEvent(userID.String(), occurTime.Value()),

		Username:  username,
		Email:     email,
		FullName:  fullName,
		AvatarKey: avatarKey,
	}
}

func (e UserCreated) Type() string {
	return "user.created"
}

func (e UserCreated) AggregateType() string {
	return "user"
}
