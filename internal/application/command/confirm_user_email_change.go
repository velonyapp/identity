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

type ConfirmUserEmailChange struct {
	UserID string
	Token  string
}

type ConfirmUserEmailChangeResult struct{}

type ConfirmUserEmailChangeHandler struct {
	userRepo     repo.User
	unitOfWork   port.UnitOfWork
	cache        port.Cache
	verifier     port.RequestVerifier
	availability *service.EmailAvailability
}

func NewConfirmUserEmailChangeHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	cache port.Cache,
	verifier port.RequestVerifier,
	availability *service.EmailAvailability,
) *ConfirmUserEmailChangeHandler {
	return &ConfirmUserEmailChangeHandler{
		userRepo:     userRepo,
		unitOfWork:   unitOfWork,
		cache:        cache,
		verifier:     verifier,
		availability: availability,
	}
}

func (h *ConfirmUserEmailChangeHandler) Execute(
	ctx context.Context,
	cmd *ConfirmUserEmailChange,
) (*ConfirmUserEmailChangeResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		user, err := h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		if user.HasEmailChangeRequest() {
			if err := h.verifier.VerifyEmailChange(*user.EmailChangeRequest(), cmd.Token); err != nil {
				return err
			}
		}

		if err := user.ConfirmEmailChange(now); err != nil {
			return err
		}

		if err := h.availability.EnsureAvailable(ctx, *user.Email()); err != nil {
			return err
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	h.cache.Delete(ctx, common.UserResultCacheKey(userID))

	return &ConfirmUserEmailChangeResult{}, nil
}
