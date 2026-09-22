package command

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrEmailChangeRequestNotFound = errors.New("email change request not found")

type ConfirmUserEmailChange struct {
	Token string
}

type ConfirmUserEmailChangeResult struct{}

type ConfirmUserEmailChangeHandler struct {
	requestRepo       repo.EmailChangeRequest
	userRepo          repo.User
	unitOfWork        port.UnitOfWork
	cache             port.Cache
	emailAvailability *service.EmailAvailability
}

func NewConfirmUserEmailChangeHandler(
	requestRepo repo.EmailChangeRequest,
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
	emailAvailability *service.EmailAvailability,
) *ConfirmUserEmailChangeHandler {
	return &ConfirmUserEmailChangeHandler{
		requestRepo:       requestRepo,
		userRepo:          userRepo,
		unitOfWork:        unitOfWork,
		cache:             cache,
		emailAvailability: emailAvailability,
	}
}

func (h *ConfirmUserEmailChangeHandler) Execute(
	ctx context.Context,
	cmd *ConfirmUserEmailChange,
) (*ConfirmUserEmailChangeResult, error) {
	token, err := vo.NewRequestToken(cmd.Token)
	if err != nil {
		return nil, err
	}

	var user *entity.User

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		request, err := h.requestRepo.FindByToken(ctx, token)
		if err != nil {
			return err
		}

		if request == nil {
			return ErrEmailChangeRequestNotFound
		}

		if err := request.Confirm(); err != nil {
			return err
		}

		if err := h.emailAvailability.EnsureAvailable(ctx, request.NewEmail); err != nil {
			return err
		}

		user, err = h.userRepo.FindByID(ctx, request.UserID)
		if err != nil {
			return err
		}

		if err := user.ChangeEmail(&request.NewEmail); err != nil {
			return err
		}

		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return h.requestRepo.Save(ctx, request)
	}); err != nil {
		return nil, err
	}

	h.cache.Set(ctx,
		common.UserResultCacheKey(user.ID),
		common.NewUserResult(user),
		common.UserResultCacheTTL,
	)

	return &ConfirmUserEmailChangeResult{}, nil
}
