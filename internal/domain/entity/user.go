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
	ErrEmailChangeNotRequested      = errors.New("email change not requested")
	ErrEmailChangeRequestExpired    = errors.New("email change request expired")
	ErrAvatarAlreadySet             = errors.New("avatar already set")
	ErrAvatarChangeNotRequested     = errors.New("avatar change not requested")
	ErrAvatarChangeRequestExpired   = errors.New("avatar change request expired")
	ErrLocalAuthStrategyNotSet      = errors.New("local authentication strategy not set")
	ErrGoogleAuthStrategyNotSet     = errors.New("google authentication strategy not set")
	ErrCannotRemoveLastAuthStrategy = errors.New("cannot remove the last authentication strategy")
)

type User struct {
	id         vo.UserID
	username   vo.Username
	email      *vo.Email
	avatarID   *vo.AvatarID
	createTime time.Time
	deleteTime *time.Time

	emailChangeRequest  *vo.EmailChangeRequest
	avatarChangeRequest *vo.AvatarChangeRequest

	localAuthStrategy *vo.LocalAuthStrategy

	domainEvents []event.DomainEvent

	version uint32
	isNew   bool
}

func NewUser(
	username vo.Username,
	now time.Time,
) *User {
	u := &User{
		id:         vo.GenerateUserID(),
		username:   username,
		createTime: now,

		version: 0,
		isNew:   true,
	}

	u.recordEvent(
		event.NewUserCreated(
			u.id,
			u.username,
			u.createTime,
			u.version,
		),
	)

	return u
}

func ReconstituteUser(
	id vo.UserID,
	username vo.Username,
	email *vo.Email,
	avatarID *vo.AvatarID,
	createTime time.Time,
	deleteTime *time.Time,
	emailChangeRequest *vo.EmailChangeRequest,
	avatarChangeRequest *vo.AvatarChangeRequest,
	localAuthStrategy *vo.LocalAuthStrategy,
	version uint32,
) *User {
	u := &User{
		id:         id,
		username:   username,
		createTime: createTime,

		version: version,
	}

	if email != nil {
		value := *email
		u.email = &value
	}
	if avatarID != nil {
		value := *avatarID
		u.avatarID = &value
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
	if localAuthStrategy != nil {
		value := *localAuthStrategy
		u.localAuthStrategy = &value
	}

	return u
}

func (u *User) ID() vo.UserID {
	return u.id
}

func (u *User) Username() vo.Username {
	return u.username
}

func (u *User) Email() *vo.Email {
	if u.email == nil {
		return nil
	}

	value := *u.email
	return &value
}

func (u *User) AvatarID() *vo.AvatarID {
	if u.avatarID == nil {
		return nil
	}

	value := *u.avatarID
	return &value
}

func (u *User) CreateTime() time.Time {
	return u.createTime
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

func (u *User) HasAvatar() bool {
	return u.avatarID != nil
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

func (u *User) AuthStrategyCount() int {
	count := 0

	if u.localAuthStrategy != nil {
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
	if newUsername.Equal(u.username) {
		return nil
	}

	u.username = newUsername

	u.recordEvent(
		event.NewUserUpdated(
			u.id,
			u.username,
			u.email,
			now,
			u.version,
		),
	)

	return nil
}

func (u *User) RequestEmailChange(newEmail vo.Email, token vo.EmailChangeRequestToken, ttl time.Duration, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if u.HasEmail() && newEmail.Equal(*u.email) {
		return ErrEmailAlreadySet
	}

	value := vo.NewEmailChangeRequest(
		newEmail,
		token,
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
		return ErrEmailChangeRequestExpired
	}
	if u.HasEmail() && u.emailChangeRequest.Email().Equal(*u.email) {
		return nil
	}

	value := u.emailChangeRequest.Email()
	newEmail := &value
	u.email = newEmail
	u.emailChangeRequest = nil

	u.recordEvent(
		event.NewUserUpdated(
			u.id,
			u.username,
			u.email,
			now,
			u.version,
		),
	)

	return nil
}

func (u *User) RemoveEmail(now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if !u.HasEmail() {
		return nil
	}

	u.email = nil

	u.recordEvent(
		event.NewUserUpdated(
			u.id,
			u.username,
			u.email,
			now,
			u.version,
		),
	)

	return nil
}

func (u *User) RequestAvatarChange(newAvatarID vo.AvatarID, ttl time.Duration, now time.Time) error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}
	if u.HasAvatar() && newAvatarID.Equal(*u.avatarID) {
		return ErrAvatarAlreadySet
	}

	value := vo.NewAvatarChangeRequest(
		newAvatarID,
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
	if u.HasAvatar() && u.avatarChangeRequest.AvatarID().Equal(*u.avatarID) {
		return nil
	}

	value := u.avatarChangeRequest.AvatarID()
	newAvatar := &value
	u.avatarID = newAvatar
	u.avatarChangeRequest = nil

	return nil
}

func (u *User) RemoveAvatar() error {
	if u.IsDeleted() {
		return ErrUserDeleted
	}

	if !u.HasAvatar() {
		return nil
	}

	u.avatarID = nil

	return nil
}

func (u *User) Delete(now time.Time) error {
	if u.IsDeleted() {
		return nil
	}

	u.deleteTime = &now

	u.recordEvent(
		event.NewUserDeleted(
			u.id,
			*u.deleteTime,
			u.version,
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

func (u *User) PullEvents() []event.DomainEvent {
	pulled := u.domainEvents
	u.domainEvents = nil
	return pulled
}

func (u *User) recordEvent(domainEvent event.DomainEvent) {
	u.domainEvents = append(u.domainEvents, domainEvent)
}

func (u *User) Version() uint32 {
	return u.version
}

func (u *User) IsNew() bool {
	return u.isNew
}

func (u *User) MarkPersisted() {
	if u.isNew {
		u.isNew = false

		return
	}

	u.version++
}
