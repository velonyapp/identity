package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RemoveUserAvatar struct {
	UserID string
}

type RemoveUserAvatarResult struct {
	User *common.User
}

type RemoveUserAvatarHandler struct {
	userRepo   repo.User
	unitOfWork port.UnitOfWork
	cache      port.Cache
}

func NewRemoveUserAvatarHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
) *RemoveUserAvatarHandler {
	return &RemoveUserAvatarHandler{
		userRepo:   userRepo,
		unitOfWork: unitOfWork,
		cache:      cache,
	}
}
func (h *RemoveUserAvatarHandler) Execute(
	ctx context.Context,
	cmd *RemoveUserAvatar,
) (*RemoveUserAvatarResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)

	var result *common.User

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		user, err := h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return common.ErrUserNotFound
		}

		if err := user.RemoveAvatar(now); err != nil {
			return err
		}

		result = &common.User{
			ID:       user.ID().Value(),
			Username: user.Username().Value(),
			FullName: user.FullName().Value(),
		}
		if user.HasEmail() {
			value := user.Email().Value()
			result.Email = &value
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserCacheKey(userID.Value()))

	return &RemoveUserAvatarResult{User: result}, nil
}
