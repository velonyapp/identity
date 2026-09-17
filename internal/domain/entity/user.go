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
	ErrAvatarKeyUserMismatch = errors.New(
		"avatar key does not belong to user",
	)
)

type User struct {
	ID         vo.UserID
	Username   vo.Username
	Email      *vo.Email
	FullName   vo.FullName
	AvatarKey  *vo.AvatarKey
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
	avatarKey *vo.AvatarKey,
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

func (u *User) ChangeUsername(newUsername vo.Username) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	oldUsername := u.Username
	now := vo.NewTimeNow()

	u.Username = newUsername
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserUsernameChanged(
			u.ID,
			oldUsername,
			newUsername,
			now,
		),
	)

	return nil
}

func (u *User) ChangeEmail(newEmail *vo.Email) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	oldEmail := u.Email
	now := vo.NewTimeNow()

	u.Email = newEmail
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserEmailChanged(
			u.ID,
			oldEmail,
			newEmail,
			now,
		),
	)

	return nil
}

func (u *User) ChangeFullName(newFullName vo.FullName) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	oldFullName := u.FullName
	now := vo.NewTimeNow()

	u.FullName = newFullName
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserFullNameChanged(
			u.ID,
			oldFullName,
			newFullName,
			now,
		),
	)

	return nil
}

func (u *User) ChangeAvatarKey(newAvatarKey *vo.AvatarKey) error {
	if u.DeleteTime != nil {
		return ErrUserDeleted
	}

	if newAvatarKey != nil && newAvatarKey.UserID() != u.ID.Value() {
		return ErrAvatarKeyUserMismatch
	}

	oldAvatarKey := u.AvatarKey
	now := vo.NewTimeNow()

	u.AvatarKey = newAvatarKey
	u.UpdateTime = now

	u.recordEvent(
		event.NewUserAvatarKeyChanged(
			u.ID,
			oldAvatarKey,
			newAvatarKey,
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
