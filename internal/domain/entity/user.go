package entity

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/event"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrUserDeleted = errors.New(
		"user is deleted",
	)
	ErrCannotRemoveLastAuthStrategy = errors.New(
		"cannot remove the last authentication strategy",
	)
)

type User struct {
	ID         vo.UserID
	Username   vo.Username
	Email      *vo.Email
	FullName   vo.FullName
	AvatarKey  *vo.StorageKey
	CreateTime vo.Time
	UpdateTime vo.Time
	DeleteTime *vo.Time

	LocalAuthStrategy  *LocalAuthStrategy
	GoogleAuthStrategy *GoogleAuthStrategy

	domainEvents []event.DomainEvent
}

func NewUser(
	username vo.Username,
	email *vo.Email,
	fullName vo.FullName,
	avatarKey *vo.StorageKey,
) *User {
	now := vo.NewTimeNow()
	userID := vo.NewUserIDRandom()

	user := &User{
		ID:         userID,
		Username:   username,
		FullName:   fullName,
		CreateTime: now,
		UpdateTime: now,
	}

	user.recordEvent(
		event.NewUserCreated(
			userID,
			username,
			email,
			fullName,
			avatarKey,
			now,
		),
	)

	return user
}

func (u *User) ChangeUsername(username vo.Username) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	now := vo.NewTimeNow()

	u.Username = username
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserUsernameChanged(
			u.ID,
			username,
			now,
		),
	)

	return nil
}

func (u *User) ChangeEmail(email *vo.Email) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	now := vo.NewTimeNow()

	u.Email = email
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserEmailChanged(
			u.ID,
			email,
			now,
		),
	)

	return nil
}

func (u *User) ChangeFullName(fullName vo.FullName) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	now := vo.NewTimeNow()

	u.FullName = fullName
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserFullNameChanged(
			u.ID,
			fullName,
			now,
		),
	)

	return nil
}

func (u *User) ChangeAvatarKey(avatarKey *vo.StorageKey) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	now := vo.NewTimeNow()

	u.AvatarKey = avatarKey
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserAvatarKeyChanged(
			u.ID,
			avatarKey,
			now,
		),
	)

	return nil
}

func (u *User) Delete() error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	now := vo.NewTimeNow()

	u.UpdateTime = now
	u.DeleteTime = &now

	u.recordEvent(
		event.NewUserDeleted(
			u.ID,
			now,
		),
	)

	return nil
}

func (u *User) ReplaceLocalAuthStrategy(localAuthStrategy *LocalAuthStrategy) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}
	if localAuthStrategy == nil && u.GoogleAuthStrategy == nil {
		return ErrCannotRemoveLastAuthStrategy
	}

	u.LocalAuthStrategy = localAuthStrategy

	return nil
}

func (u *User) ReplaceGoogleAuthStrategy(googleAuthStrategy *GoogleAuthStrategy) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}
	if googleAuthStrategy == nil && u.LocalAuthStrategy == nil {
		return ErrCannotRemoveLastAuthStrategy
	}

	u.GoogleAuthStrategy = googleAuthStrategy

	return nil
}

func (u *User) PullEvents() []event.DomainEvent {
	pulled := u.domainEvents
	u.domainEvents = nil
	return pulled
}

func (u *User) recordEvent(domainEvent event.DomainEvent) {
	u.domainEvents = append(u.domainEvents, domainEvent)
}
