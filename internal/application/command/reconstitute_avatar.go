package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type ReconstituteAvatar struct {
	AvatarID string
	Key      string
}

type ReconstituteAvatarResult struct {
}

type ReconstituteAvatarHandler struct {
	userRepo   repo.User
	avatarRepo repo.Avatar
	unitOfWork port.UnitOfWork
}

func NewReconstituteAvatarHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
) *ReconstituteAvatarHandler {
	return &ReconstituteAvatarHandler{
		userRepo:   userRepo,
		avatarRepo: avatarRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *ReconstituteAvatarHandler) Execute(
	ctx context.Context,
	cmd *ReconstituteAvatar,
) (*ReconstituteAvatarResult, error) {
	avatarID := vo.NewAvatarID(cmd.AvatarID)
	key, err := vo.NewAvatarKey(cmd.Key)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByRequestedAvatarIDChange(ctx, avatarID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	avatar := entity.ReconstituteAvatar(
		avatarID,
		user.ID(),
		key,
		nil,
	)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.avatarRepo.Save(ctx, avatar)
	}); err != nil {
		return nil, err
	}

	return &ReconstituteAvatarResult{}, nil
}
