package common

import (
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/domain/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        string
	Username  string
	FullName  string
	Email     *string
	AvatarKey *string
}

func NewUserFromEntity(user *entity.User) *User {
	result := &User{
		ID:       user.ID().Value(),
		Username: user.Username().Value(),
		FullName: user.FullName().Value(),
	}
	if user.HasEmail() {
		value := user.Email().Value()
		result.Email = &value
	}
	if user.HasAvatar() {
		value := user.AvatarKey().String()
		result.AvatarKey = &value
	}

	return result
}

func UserCacheKey(userID string) string {
	return "user:" + userID
}

const UserCacheTTL = 1 * time.Minute
