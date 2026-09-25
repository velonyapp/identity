package common

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type UserResult struct {
	ID        string
	Username  string
	FullName  string
	Email     *string
	AvatarKey *string
}

func NewUserResult(user *entity.User) *UserResult {
	result := &UserResult{
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

func UserResultCacheKey(userID vo.UserID) string {
	return "user:" + userID.Value()
}

const UserResultCacheTTL = 15 * time.Minute
