package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type DeleteAvatar struct {
	AvatarID string
}

type DeleteAvatarResult struct {
}

type DeleteAvatarHandler struct {
	avatarRepo repo.Avatar
	unitOfWork port.UnitOfWork
}

func NewDeleteAvatarHandler(
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
) *DeleteAvatarHandler {
	return &DeleteAvatarHandler{
		avatarRepo: avatarRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *DeleteAvatarHandler) Execute(
	ctx context.Context,
	cmd *DeleteAvatar,
) (*DeleteAvatarResult, error) {
	now := time.Now()

	avatarID := vo.NewAvatarID(cmd.AvatarID)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		avatar, err := h.avatarRepo.FindByID(ctx, avatarID)
		if err != nil {
			return err
		}

		avatar.Delete(now)

		if err := h.avatarRepo.Save(ctx, avatar); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &DeleteAvatarResult{}, nil
}
