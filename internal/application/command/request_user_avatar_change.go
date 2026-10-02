package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RequestUserAvatarChange struct {
	UserID string
}

type RequestUserAvatarChangeResult struct {
	UploadURL string
}

type RequestUserAvatarChangeHandler struct {
	c            *conf.Security
	userRepo     repo.User
	unitOfWork   port.UnitOfWork
	assetService port.AssetService
}

func NewRequestUserAvatarChangeHandler(
	c *conf.Security,
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	assetService port.AssetService,
) *RequestUserAvatarChangeHandler {
	return &RequestUserAvatarChangeHandler{
		c:            c,
		userRepo:     userRepo,
		unitOfWork:   unitOfWork,
		assetService: assetService,
	}
}

func (h *RequestUserAvatarChangeHandler) Execute(
	ctx context.Context,
	cmd *RequestUserAvatarChange,
) (*RequestUserAvatarChangeResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.IsDeleted() {
		return nil, common.ErrUserNotFound
	}

	avatarID, err := h.assetService.CreateAvatar(ctx, userID)
	if err != nil {
		return nil, err
	}
	uploadURL, err := h.assetService.PresignAvatar(ctx, avatarID, h.c.AvatarChangeRequest.Ttl.AsDuration())
	if err != nil {
		return nil, err
	}

	if err := user.RequestAvatarChange(avatarID, h.c.AvatarChangeRequest.Ttl.AsDuration(), now); err != nil {
		return nil, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	return &RequestUserAvatarChangeResult{
		UploadURL: uploadURL,
	}, nil
}
