package common

import (
	"time"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type UserResult struct {
	ID        string
	Username  string
	Email     *string
	FullName  string
	AvatarKey *string
}

func NewUserResult(user *entity.User) *UserResult {
	result := &UserResult{
		ID:       user.ID.String(),
		Username: user.Username.Value(),
		FullName: user.FullName.Value(),
	}
	if user.Email != nil {
		value := user.Email.Value()
		result.Email = &value
	}
	if user.AvatarKey != nil {
		value := user.AvatarKey.String()
		result.AvatarKey = &value
	}

	return result
}

func UserResultCacheKey(userID vo.UserID) string {
	return "user:" + userID.Value()
}

const UserResultCacheTTL = 15 * time.Minute
