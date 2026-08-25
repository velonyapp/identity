package event

import "github.com/velony-app/identity/internal/domain/vo"

type UserCreated struct {
	BaseDomainEvent

	Username   vo.Username
	Email      *vo.Email
	FullName   vo.FullName
	AvatarKey  *vo.StorageKey
	CreateTime vo.Time
}

func NewUserCreated(
	userID vo.UserID,
	username vo.Username,
	email *vo.Email,
	fullName vo.FullName,
	avatarKey *vo.StorageKey,
	createTime vo.Time,
) UserCreated {
	return UserCreated{
		BaseDomainEvent: NewBaseDomainEvent(userID.String()),

		Username:   username,
		Email:      email,
		FullName:   fullName,
		AvatarKey:  avatarKey,
		CreateTime: createTime,
	}
}

func (e UserCreated) Type() string {
	return "user.created"
}
