package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RequestUserEmailChange struct {
	UserID string
	Email  string
}

type RequestUserEmailChangeResult struct{}

type RequestUserEmailChangeHandler struct {
	c             *conf.Security
	userRepo      repo.User
	unitOfWork    port.UnitOfWork
	tokenProvider port.AuthTokenProvider
	availability  *service.EmailAvailability
}

func NewRequestUserEmailChangeHandler(
	c *conf.Security,
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	tokenProvider port.AuthTokenProvider,
	availability *service.EmailAvailability,
) *RequestUserEmailChangeHandler {
	return &RequestUserEmailChangeHandler{
		c:             c,
		userRepo:      userRepo,
		unitOfWork:    unitOfWork,
		tokenProvider: tokenProvider,
		availability:  availability,
	}
}

func (h *RequestUserEmailChangeHandler) Execute(
	ctx context.Context,
	cmd *RequestUserEmailChange,
) (*RequestUserEmailChangeResult, error) {
	now := time.Now()

	userID := vo.NewUserID(cmd.UserID)
	newEmail, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.availability.EnsureAvailable(ctx, newEmail); err != nil {
			return err
		}

		user, err := h.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}

		if err := user.RequestEmailChange(newEmail, h.c.EmailChangeRequestToken.Ttl.AsDuration(), now); err != nil {
			return err
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return nil, err
	}

	return &RequestUserEmailChangeResult{}, nil
}
