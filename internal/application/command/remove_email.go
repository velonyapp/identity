package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RemoveEmail struct {
	UserID string
}

type RemoveEmailResult struct {
	User common.UserResult
}

func (RemoveEmail) resultType() RemoveEmailResult {
	return RemoveEmailResult{}
}

type RemoveEmailHandler Handler[RemoveEmail, RemoveEmailResult]

type removeEmailHandler struct {
	userRepo   repo.User
	unitOfWork port.UnitOfWork
}

func NewRemoveEmailHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
) RemoveEmailHandler {
	return &removeEmailHandler{
		userRepo:   userRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *removeEmailHandler) Handle(
	ctx context.Context,
	cmd RemoveEmail,
) (RemoveEmailResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return RemoveEmailResult{}, err
	}
	if user == nil {
		return RemoveEmailResult{}, common.ErrUserNotFound
	}

	if err := user.RemoveEmail(now); err != nil {
		return RemoveEmailResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return RemoveEmailResult{}, err
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

	return RemoveEmailResult{
		User: userResult,
	}, nil
}
