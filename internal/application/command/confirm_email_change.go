package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type ConfirmEmailChange struct {
	UserID string
	Token  string
}

type ConfirmEmailChangeResult struct {
	User common.UserResult
}

func (ConfirmEmailChange) resultType() ConfirmEmailChangeResult {
	return ConfirmEmailChangeResult{}
}

type ConfirmEmailChangeHandler Handler[ConfirmEmailChange, ConfirmEmailChangeResult]

type confirmEmailChangeHandler struct {
	userRepo            repo.User
	avatarRepo          repo.Avatar
	unitOfWork          port.UnitOfWork
	requestTokenManager port.EmailChangeRequestTokenManager
}

func NewConfirmEmailChangeHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	requestTokenManager port.EmailChangeRequestTokenManager,
) ConfirmEmailChangeHandler {
	return &confirmEmailChangeHandler{
		userRepo:            userRepo,
		avatarRepo:          avatarRepo,
		unitOfWork:          unitOfWork,
		requestTokenManager: requestTokenManager,
	}
}

func (h *confirmEmailChangeHandler) Handle(
	ctx context.Context,
	cmd ConfirmEmailChange,
) (ConfirmEmailChangeResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ConfirmEmailChangeResult{}, err
	}
	if user == nil {
		return ConfirmEmailChangeResult{}, common.ErrUserNotFound
	}

	if user.HasEmailChangeRequest() {
		requestToken, _ := vo.NewEmailChangeRequestToken(cmd.Token)

		if err := h.requestTokenManager.Verify(requestToken, user.EmailChangeRequest().Token()); err != nil {
			return ConfirmEmailChangeResult{}, err
		}
	}

	if err := user.ConfirmEmailChange(now); err != nil {
		return ConfirmEmailChangeResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return ConfirmEmailChangeResult{}, err
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
			return ConfirmEmailChangeResult{}, err
		}

		value := avatar.Key().String()
		userResult.AvatarKey = &value
	}

	return ConfirmEmailChangeResult{
		User: userResult,
	}, nil
}
