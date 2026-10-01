package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type UpdateUser struct {
	UserID   string
	Username *string
	FullName *string
}

type UpdateUserResult struct {
	User *common.User
}

type UpdateUserHandler struct {
	userRepo             repo.User
	unitOfWork           port.UnitOfWork
	cache                port.Cache
	usernameAvailability *service.UsernameAvailability
}

func NewUpdateUserHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
	usernameAvailability *service.UsernameAvailability,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepo:             userRepo,
		unitOfWork:           unitOfWork,
		cache:                cache,
		usernameAvailability: usernameAvailability,
	}
}
func (h *UpdateUserHandler) Execute(
	ctx context.Context,
	cmd *UpdateUser,
) (*UpdateUserResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)

	var user *entity.User

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		user, err = h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if user == nil {
			return common.ErrUserNotFound
		}

		if cmd.Username != nil {
			newUsername, err := vo.NewUsername(*cmd.Username)
			if err != nil {
				return err
			}

			if newUsername != user.Username() {
				if err := h.usernameAvailability.EnsureAvailable(ctx, newUsername); err != nil {
					return err
				}
			}

			if err := user.ChangeUsername(newUsername, now); err != nil {
				return err
			}
		}
		if cmd.FullName != nil {
			fullName, err := vo.NewFullName(*cmd.FullName)
			if err != nil {
				return err
			}

			if err := user.ChangeFullName(fullName, now); err != nil {
				return err
			}
		}

		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserCacheKey(userID.Value()))

	return &UpdateUserResult{User: common.NewUserFromEntity(user)}, nil
}
