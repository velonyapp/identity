package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type ConfirmUserAvatarChange struct {
	AvatarID string
}

type ConfirmUserAvatarChangeResult struct{}

type ConfirmUserAvatarChangeHandler struct {
	userRepo   repo.User
	avatarRepo repo.Avatar
	unitOfWork port.UnitOfWork
	cache      port.Cache
}

func NewConfirmUserAvatarChangeHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
	verifier port.RequestVerifier,
) *ConfirmUserAvatarChangeHandler {
	return &ConfirmUserAvatarChangeHandler{
		userRepo:   userRepo,
		avatarRepo: avatarRepo,
		unitOfWork: unitOfWork,
		cache:      cache,
	}
}

func (h *ConfirmUserAvatarChangeHandler) Execute(
	ctx context.Context,
	cmd *ConfirmUserAvatarChange,
) (*ConfirmUserAvatarChangeResult, error) {
	now := time.Now()

	avatarID := vo.NewAvatarID(cmd.AvatarID)

	avatar, err := h.avatarRepo.FindByID(ctx, avatarID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, avatar.UserID())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	if err := user.ConfirmAvatarChange(now); err != nil {
		return nil, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserCacheKey(user.ID().Value()))

	return &ConfirmUserAvatarChangeResult{}, nil
}
