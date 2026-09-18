package integrationevent

import "time"

type UserCreated struct {
	BaseIntegrationEvent

	Username   string
	Email      *string
	FullName   string
	AvatarKey  *string
	CreateTime time.Time
}

func NewUserCreated(
	userID string,
	username string,
	email *string,
	fullName string,
	avatarKey *string,
	createTime time.Time,
) UserCreated {
	return UserCreated{
		BaseIntegrationEvent: NewBaseIntegrationEvent(userID),

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

func (e UserCreated) AggregateType() string {
	return "user"
}
