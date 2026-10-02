package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
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
	avatarRepo           repo.Avatar
	unitOfWork           port.UnitOfWork
	cache                port.Cache
	usernameAvailability *service.UsernameAvailability
}

func NewUpdateUserHandler(
	userRepo repo.User,
	avatarRepo repo.Avatar,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
	usernameAvailability *service.UsernameAvailability,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepo:             userRepo,
		avatarRepo:           avatarRepo,
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

	var result *common.User

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		user, err := h.userRepo.FindByID(ctx, userID)
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

		result = &common.User{
			ID:       user.ID().Value(),
			Username: user.Username().Value(),
			FullName: user.FullName().Value(),
		}
		if user.HasEmail() {
			value := user.Email().Value()
			result.Email = &value
		}
		if user.HasAvatar() {
			avatar, err := h.avatarRepo.FindByID(ctx, *user.AvatarID())
			if err != nil {
				return err
			}

			value := avatar.Key().Value()
			result.AvatarKey = &value
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserCacheKey(userID.Value()))

	return &UpdateUserResult{User: result}, nil
}
