package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type UpdateUser struct {
	UserID    string
	Username  *string
	FullName  *string
	AvatarKey **string
}

type UpdateUserResult struct {
	User *common.UserResult
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

		changed := false

		if cmd.Username != nil {
			username, err := vo.NewUsername(*cmd.Username)
			if err != nil {
				return err
			}

			if username != user.Username {
				if err := h.usernameAvailability.EnsureAvailable(ctx, username); err != nil {
					return err
				}

				if err := user.ChangeUsername(username); err != nil {
					return err
				}

				changed = true
			}
		}
		if cmd.FullName != nil {
			fullName, err := vo.NewFullName(*cmd.FullName)
			if err != nil {
				return err
			}

			if fullName != user.FullName {
				if err := user.ChangeFullName(fullName); err != nil {
					return err
				}

				changed = true
			}
		}
		if cmd.AvatarKey != nil {
			var avatarKey *vo.AvatarKey

			if *cmd.AvatarKey != nil {
				value, err := vo.NewAvatarKey(**cmd.AvatarKey)
				if err != nil {
					return err
				}

				avatarKey = &value
			}

			if avatarKey != user.AvatarKey {
				if err := user.ChangeAvatarKey(avatarKey); err != nil {
					return err
				}

				changed = true
			}
		}

		if changed {
			if err := h.userRepo.Save(ctx, user); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	userResult := common.NewUserResult(user)

	h.cache.Set(ctx,
		common.UserResultCacheKey(user.ID),
		userResult,
		common.UserResultCacheTTL,
	)

	return &UpdateUserResult{User: userResult}, nil
}
