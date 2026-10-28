package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/job"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RemoveAvatar struct {
	UserID string
}

type RemoveAvatarResult struct {
	User common.UserResult
}

func (RemoveAvatar) resultType() RemoveAvatarResult {
	return RemoveAvatarResult{}
}

type RemoveAvatarHandler Handler[RemoveAvatar, RemoveAvatarResult]

type removeAvatarHandler struct {
	userRepo     repo.User
	avatarRepo   repo.Avatar
	unitOfWork   port.UnitOfWork
	jobPublisher port.JobPublisher
}

func NewRemoveAvatarHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	jobPublisher port.JobPublisher,
) RemoveAvatarHandler {
	return &removeAvatarHandler{
		userRepo:     userRepo,
		avatarRepo:   avatarRepo,
		unitOfWork:   unitOfWork,
		jobPublisher: jobPublisher,
	}
}

func (h *removeAvatarHandler) Handle(
	ctx context.Context,
	cmd RemoveAvatar,
) (RemoveAvatarResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return RemoveAvatarResult{}, err
	}
	if user == nil {
		return RemoveAvatarResult{}, common.ErrUserNotFound
	}

	var oldAvatarID *vo.AvatarID
	if user.HasAvatar() {
		oldAvatarID = user.AvatarID()
	}

	if err := user.RemoveAvatar(); err != nil {
		return RemoveAvatarResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if oldAvatarID != nil {
			if err := h.avatarRepo.Delete(ctx, *oldAvatarID); err != nil {
				return err
			}

			if err := h.jobPublisher.Publish(ctx, job.NewRemoveAvatar(oldAvatarID.String(), now)); err != nil {
				return err
			}
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return RemoveAvatarResult{}, err
	}

	userResult := common.UserResult{
		ID:         user.ID().String(),
		Username:   user.Username().String(),
		CreateTime: user.CreateTime(),
	}
	if user.HasEmail() {
		value := user.Email().String()
		userResult.Email = &value
	}

	return RemoveAvatarResult{
		User: userResult,
	}, nil
}
