package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type DeleteUser struct {
	UserID string
}

type DeleteUserResult struct {
}

type DeleteUserHandler struct {
	userRepo   repo.User
	unitOfWork port.UnitOfWork
	cache      port.Cache
}

func NewDeleteUserHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,

) *DeleteUserHandler {
	return &DeleteUserHandler{
		userRepo:   userRepo,
		unitOfWork: unitOfWork,
		cache:      cache,
	}
}

func (h *DeleteUserHandler) Execute(
	ctx context.Context,
	cmd *DeleteUser,
) (*DeleteUserResult, error) {
	userID := vo.NewUserID(cmd.UserID)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		user, err := h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return common.ErrUserNotFound
		}

		if err := user.Delete(); err != nil {
			return err
		}

		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserResultCacheKey(userID))

	return &DeleteUserResult{}, nil
}
