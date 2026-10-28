package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type ChangeUsername struct {
	UserID   string
	Username string
}

type ChangeUsernameResult struct {
	User common.UserResult
}

func (ChangeUsername) resultType() ChangeUsernameResult {
	return ChangeUsernameResult{}
}

type ChangeUsernameHandler Handler[ChangeUsername, ChangeUsernameResult]

type changeUsernameHandler struct {
	userRepo   repo.User
	avatarRepo repo.Avatar
	unitOfWork port.UnitOfWork
}

func NewChangeUsernameHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
) ChangeUsernameHandler {
	return &changeUsernameHandler{
		userRepo:   userRepo,
		avatarRepo: avatarRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *changeUsernameHandler) Handle(
	ctx context.Context,
	cmd ChangeUsername,
) (ChangeUsernameResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ChangeUsernameResult{}, err
	}
	if user == nil {
		return ChangeUsernameResult{}, common.ErrUserNotFound
	}

	newUsername, err := vo.NewUsername(cmd.Username)
	if err != nil {
		return ChangeUsernameResult{}, err
	}

	if err := user.ChangeUsername(newUsername, now); err != nil {
		return ChangeUsernameResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return ChangeUsernameResult{}, err
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
			return ChangeUsernameResult{}, err
		}

		value := avatar.Key().String()
		userResult.AvatarKey = &value
	}

	return ChangeUsernameResult{
		User: userResult,
	}, nil
}
