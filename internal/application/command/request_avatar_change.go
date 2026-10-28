package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RequestAvatarChange struct {
	UserID string
}

type RequestAvatarChangeResult struct {
	UploadURL string
}

func (RequestAvatarChange) resultType() RequestAvatarChangeResult {
	return RequestAvatarChangeResult{}
}

type RequestAvatarChangeHandler Handler[RequestAvatarChange, RequestAvatarChangeResult]

type requestAvatarChangeHandler struct {
	userRepo     repo.User
	avatarRepo   repo.Avatar
	unitOfWork   port.UnitOfWork
	jobPublisher port.JobPublisher
	assetService port.AssetService
}

func NewRequestAvatarChangeHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	jobPublisher port.JobPublisher,
	assetService port.AssetService,
) RequestAvatarChangeHandler {
	return &requestAvatarChangeHandler{
		userRepo:     userRepo,
		avatarRepo:   avatarRepo,
		unitOfWork:   unitOfWork,
		jobPublisher: jobPublisher,
		assetService: assetService,
	}
}

func (h *requestAvatarChangeHandler) Handle(
	ctx context.Context,
	cmd RequestAvatarChange,
) (RequestAvatarChangeResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return RequestAvatarChangeResult{}, err
	}
	if user == nil {
		return RequestAvatarChangeResult{}, common.ErrUserNotFound
	}

	avatar, uploadURL, err := h.assetService.CreateAvatar(ctx, userID)
	if err != nil {
		return RequestAvatarChangeResult{}, err
	}

	if err := user.RequestAvatarChange(avatar.ID(), common.AvatarChangeRequestTTL, now); err != nil {
		return RequestAvatarChangeResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.avatarRepo.Create(ctx, avatar); err != nil {
			return err
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return RequestAvatarChangeResult{}, err
	}

	return RequestAvatarChangeResult{
		UploadURL: uploadURL,
	}, nil
}
