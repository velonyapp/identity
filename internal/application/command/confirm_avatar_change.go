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

type ConfirmAvatarChange struct {
	UserID string
}

type ConfirmAvatarChangeResult struct {
	User common.UserResult
}

func (ConfirmAvatarChange) resultType() ConfirmAvatarChangeResult {
	return ConfirmAvatarChangeResult{}
}

type ConfirmAvatarChangeHandler Handler[ConfirmAvatarChange, ConfirmAvatarChangeResult]

type confirmAvatarChangeHandler struct {
	userRepo     repo.User
	avatarRepo   repo.Avatar
	unitOfWork   port.UnitOfWork
	jobPublisher port.JobPublisher
}

func NewConfirmAvatarChangeHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	jobPublisher port.JobPublisher,
) ConfirmAvatarChangeHandler {
	return &confirmAvatarChangeHandler{
		userRepo:     userRepo,
		avatarRepo:   avatarRepo,
		unitOfWork:   unitOfWork,
		jobPublisher: jobPublisher,
	}
}

func (h *confirmAvatarChangeHandler) Handle(
	ctx context.Context,
	cmd ConfirmAvatarChange,
) (ConfirmAvatarChangeResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ConfirmAvatarChangeResult{}, err
	}
	if user == nil {
		return ConfirmAvatarChangeResult{}, common.ErrUserNotFound
	}

	var oldAvatarID *vo.AvatarID
	if user.HasAvatar() {
		oldAvatarID = user.AvatarID()
	}

	if err := user.ConfirmAvatarChange(now); err != nil {
		return ConfirmAvatarChangeResult{}, err
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
		return ConfirmAvatarChangeResult{}, err
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
	if user.HasAvatar() {
		avatar, err := h.avatarRepo.FindByID(ctx, *user.AvatarID())
		if err != nil {
			return ConfirmAvatarChangeResult{}, err
		}

		value := avatar.Key().String()
		userResult.AvatarKey = &value
	}

	return ConfirmAvatarChangeResult{
		User: userResult,
	}, nil
}
