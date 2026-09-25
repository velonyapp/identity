package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/domain/event"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrUserDeleted                  = errors.New("user is deleted")
	ErrEmailAlreadySet              = errors.New("email already set")
	ErrEmailNotSet                  = errors.New("email not set")
	ErrEmailChangeNotRequested      = errors.New("email change not requested")
	ErrPendingEmailChangeExpired    = errors.New("email change request expired")
	ErrAvatarKeyMismatch            = errors.New("avatar does not belong to user")
	ErrAvatarNotSet                 = errors.New("avatar not set")
	ErrAvatarChangeNotRequested     = errors.New("avatar change not requested")
	ErrAvatarChangeRequestExpired   = errors.New("avatar request expired")
	ErrLocalAuthStrategyNotSet      = errors.New("local authentication strategy not set")
	ErrGoogleAuthStrategyNotSet     = errors.New("google authentication strategy not set")
	ErrCannotRemoveLastAuthStrategy = errors.New("cannot remove the last authentication strategy")
)

type User struct {
	id         vo.UserID
	username   vo.Username
	fullName   vo.FullName
	email      *vo.Email
	avatarKey  *vo.AvatarKey
	createTime time.Time
	updateTime time.Time
	deleteTime *time.Time

	emailChangeRequest  *vo.EmailChangeRequest
	avatarChangeRequest *vo.AvatarChangeRequest

	localAuthStrategy  *vo.LocalAuthStrategy
	googleAuthStrategy *vo.GoogleAuthStrategy

	domainEvents []event.DomainEvent
}

func NewUser(
	username vo.Username,
	fullName vo.FullName,
	now time.Time,
) *User {
	userID := vo.NewUserIDRandom()

	user := &User{
		id:         userID,
		username:   username,
		fullName:   fullName,
		createTime: now,
		updateTime: now,
	}

	user.recordEvent(
		event.NewUserCreated(
			userID,
			username,
			fullName,
			now,
		),
	)

	return user
}

func ReconstituteUser(
	id vo.UserID,
	username vo.Username,
	fullName vo.FullName,
	email *vo.Email,
	avatarKey *vo.AvatarKey,
	createTime time.Time,
	updateTime time.Time,
	deleteTime *time.Time,
	emailChangeRequest *vo.EmailChangeRequest,
	avatarChangeRequest *vo.AvatarChangeRequest,
	localAuthStrategy *vo.LocalAuthStrategy,
	googleAuthStrategy *vo.GoogleAuthStrategy,
) *User {
	u := &User{
		id:                 id,
		fullName:           fullName,
		username:           username,
		createTime:         createTime,
		updateTime:         updateTime,
		localAuthStrategy:  localAuthStrategy,
		googleAuthStrategy: googleAuthStrategy,
	}

	if email != nil {
		value := *email
		u.email = &value
	}
	if avatarKey != nil {
		value := *avatarKey
		u.avatarKey = &value
	}
	if deleteTime != nil {
		value := *deleteTime
		u.deleteTime = &value
	}
	if emailChangeRequest != nil {
		value := *emailChangeRequest
		u.emailChangeRequest = &value
	}
	if avatarChangeRequest != nil {
		value := *avatarChangeRequest
		u.avatarChangeRequest = &value
	}

	return u
}

func (u *User) ID() vo.UserID {
	return u.id
}

func (u *User) Username() vo.Username {
	return u.username
}

func (u *User) FullName() vo.FullName {
	return u.fullName
}

func (u *User) Email() *vo.Email {
	if u.email == nil {
		return nil
	}
	value := *u.email
	return &value
}

func (u *User) AvatarKey() *vo.AvatarKey {
	if u.avatarKey == nil {
		return nil
	}
	value := *u.avatarKey
	return &value
}

func (u *User) CreateTime() time.Time {
	return u.createTime
}

func (u *User) UpdateTime() time.Time {
	return u.updateTime
}

func (u *User) DeleteTime() *time.Time {
	if u.deleteTime == nil {
		return nil
	}
	value := *u.deleteTime
	return &value
}

func (u *User) EmailChangeRequest() *vo.EmailChangeRequest {
	if u.emailChangeRequest == nil {
		return nil
	}

	value := *u.emailChangeRequest
	return &value
}

func (u *User) AvatarChangeRequest() *vo.AvatarChangeRequest {
	if u.avatarChangeRequest == nil {
		return nil
	}

	value := *u.avatarChangeRequest
	return &value
}

func (u *User) LocalAuthStrategy() *vo.LocalAuthStrategy {
	if u.localAuthStrategy == nil {
		return nil
	}
	value := *u.localAuthStrategy
	return &value
}

func (u *User) GoogleAuthStrategy() *vo.GoogleAuthStrategy {
	if u.googleAuthStrategy == nil {
		return nil
	}
	value := *u.googleAuthStrategy
	return &value
}

func (u *User) HasAvatar() bool {
	return u.avatarKey != nil
}

func (u *User) HasEmail() bool {
	return u.email != nil
}

func (u *User) HasEmailChangeRequest() bool {
	return u.emailChangeRequest != nil
}

func (u *User) HasAvatarChangeRequest() bool {
	return u.avatarChangeRequest != nil
}

func (u *User) HasLocalAuthStrategy() bool {
	return u.localAuthStrategy != nil
}

func (u *User) HasGoogleAuthStrategy() bool {
	return u.googleAuthStrategy != nil
}

func (u *User) AuthStrategyCount() int {
	count := 0

	if u.localAuthStrategy != nil {
		count++
	}
	if u.googleAuthStrategy != nil {
		count++
	}

	return count
}

func (u *User) IsDeleted() bool {
	return u.deleteTime != nil
}

func (u *User) ChangeUsername(newUsername vo.Username, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	oldUsername := u.username

	u.username = newUsername
	u.updateTime = now

	u.recordEvent(
		event.NewUserUsernameChanged(
			u.id,
			oldUsername,
			newUsername,
			now,
		),
	)

	return nil
}

func (u *User) ChangeFullName(newFullName vo.FullName, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	oldFullName := u.fullName

	u.fullName = newFullName
	u.updateTime = now

	u.recordEvent(
		event.NewUserFullNameChanged(
			u.id,
			oldFullName,
			newFullName,
			now,
		),
	)

	return nil
}

func (u *User) RequestEmailChange(newEmail vo.Email, ttl time.Duration, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if u.HasEmail() && newEmail == *u.email {
		return ErrEmailAlreadySet
	}

	value := vo.NewEmailChangeRequest(
		newEmail,
		now,
		now.Add(ttl),
	)
	u.emailChangeRequest = &value

	return nil
}

func (u *User) ConfirmEmailChange(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasEmailChangeRequest() {
		return ErrEmailChangeNotRequested
	}
	if !now.Before(u.emailChangeRequest.ExpireTime()) {
		return ErrPendingEmailChangeExpired
	}

	var oldEmail *vo.Email
	if u.HasEmail() {
		oldEmail = u.email
	}

	value := u.emailChangeRequest.Value()
	newEmail := &value
	u.email = newEmail
	u.emailChangeRequest = nil
	u.updateTime = now

	u.recordEvent(
		event.NewUserEmailChanged(
			u.id,
			oldEmail,
			newEmail,
			now,
		),
	)

	return nil
}

func (u *User) RemoveEmail(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasEmail() {
		return ErrEmailNotSet
	}

	value := *u.email
	oldEmail := &value

	u.email = nil
	u.updateTime = now

	u.recordEvent(
		event.NewUserEmailChanged(
			u.id,
			oldEmail,
			nil,
			now,
		),
	)

	return nil
}

func (u *User) RequestAvatarChange(newAvatarKey vo.AvatarKey, ttl time.Duration, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if newAvatarKey.UserID() != u.id.Value() {
		return ErrAvatarKeyMismatch
	}

	value := vo.NewAvatarChangeRequest(
		newAvatarKey,
		now,
		now.Add(ttl),
	)
	u.avatarChangeRequest = &value

	return nil
}

func (u *User) ConfirmAvatarChange(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasAvatarChangeRequest() {
		return ErrAvatarChangeNotRequested
	}
	if !now.Before(u.avatarChangeRequest.ExpireTime()) {
		return ErrAvatarChangeRequestExpired
	}

	var oldAvatarKey *vo.AvatarKey
	if u.HasAvatar() {
		value := *u.avatarKey
		oldAvatarKey = &value
	}

	value := u.avatarChangeRequest.Value()
	newAvatarKey := &value

	u.avatarKey = newAvatarKey
	u.avatarChangeRequest = nil
	u.updateTime = now

	u.recordEvent(
		event.NewUserAvatarChanged(
			u.id,
			oldAvatarKey,
			newAvatarKey,
			now,
		),
	)

	return nil
}

func (u *User) RemoveAvatar(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasAvatar() {
		return ErrAvatarNotSet
	}

	value := *u.avatarKey
	oldAvatarKey := &value

	u.avatarKey = nil
	u.updateTime = now

	u.recordEvent(
		event.NewUserAvatarChanged(
			u.id,
			oldAvatarKey,
			nil,
			now,
		),
	)

	return nil
}

func (u *User) Delete(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	u.updateTime = now
	u.deleteTime = &now

	u.recordEvent(
		event.NewUserDeleted(
			u.id,
			now,
		),
	)

	return nil
}

func (u *User) ChangeLocalAuthStrategy(localAuthStrategy vo.LocalAuthStrategy) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	value := localAuthStrategy
	u.localAuthStrategy = &value

	return nil
}

func (u *User) RemoveLocalAuthStrategy() error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasLocalAuthStrategy() {
		return ErrLocalAuthStrategyNotSet
	}
	if u.AuthStrategyCount() == 1 {
		return ErrCannotRemoveLastAuthStrategy
	}

	u.localAuthStrategy = nil

	return nil
}

func (u *User) ChangeGoogleAuthStrategy(googleAuthStrategy vo.GoogleAuthStrategy) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	value := googleAuthStrategy
	u.googleAuthStrategy = &value

	return nil
}

func (u *User) RemoveGoogleAuthStrategy() error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasGoogleAuthStrategy() {
		return ErrGoogleAuthStrategyNotSet
	}
	if u.AuthStrategyCount() == 1 {
		return ErrCannotRemoveLastAuthStrategy
	}

	u.googleAuthStrategy = nil

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
