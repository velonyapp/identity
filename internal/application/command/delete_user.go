package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type DeleteUser struct {
	UserID string
}

type DeleteUserResult struct {
}

func (DeleteUser) resultType() DeleteUserResult {
	return DeleteUserResult{}
}

type DeleteUserHandler Handler[DeleteUser, DeleteUserResult]

type deleteUserHandler struct {
	userRepo   repo.User
	unitOfWork port.UnitOfWork
}

func NewDeleteUserHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
) DeleteUserHandler {
	return &deleteUserHandler{
		userRepo:   userRepo,
		unitOfWork: unitOfWork,
	}
}

func (h *deleteUserHandler) Handle(
	ctx context.Context,
	cmd DeleteUser,
) (DeleteUserResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return DeleteUserResult{}, err
	}
	if user == nil {
		return DeleteUserResult{}, common.ErrUserNotFound
	}

	if err := user.Delete(now); err != nil {
		return DeleteUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return DeleteUserResult{}, err
	}

	return DeleteUserResult{}, nil
}
