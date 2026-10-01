package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type ChangeUserAvatar struct {
	UserID    string
	AvatarKey *string
}

type ChangeUserAvatarResult struct {
	User *common.User
}

type ChangeUserAvatarHandler struct {
	userRepo   repo.User
	unitOfWork port.UnitOfWork
	cache      port.Cache
}

func NewChangeUserAvatarHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
) *ChangeUserAvatarHandler {
	return &ChangeUserAvatarHandler{
		userRepo:   userRepo,
		unitOfWork: unitOfWork,
		cache:      cache,
	}
}
func (h *ChangeUserAvatarHandler) Execute(
	ctx context.Context,
	cmd *ChangeUserAvatar,
) (*ChangeUserAvatarResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)

	var user *entity.User

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		user, err = h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return common.ErrUserNotFound
		}

		if cmd.AvatarKey != nil {
			newAvatarKey, err := vo.NewAvatarKey(*cmd.AvatarKey)
			if err != nil {
				return err
			}

			user.ChangeAvatar(newAvatarKey, now)
		} else {
			user.RemoveAvatar(now)
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserCacheKey(userID.Value()))

	return &ChangeUserAvatarResult{User: common.NewUserFromEntity(user)}, nil
}
