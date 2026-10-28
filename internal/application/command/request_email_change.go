package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/job"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RequestEmailChange struct {
	UserID string
	Email  string
}

type RequestEmailChangeResult struct {
}

func (RequestEmailChange) resultType() RequestEmailChangeResult {
	return RequestEmailChangeResult{}
}

type RequestEmailChangeHandler Handler[RequestEmailChange, RequestEmailChangeResult]

type requestEmailChangeHandler struct {
	userRepo            repo.User
	unitOfWork          port.UnitOfWork
	jobPublisher        port.JobPublisher
	requestTokenManager port.EmailChangeRequestTokenManager
}

func NewRequestEmailChangeHandler(
	userRepo repo.User,
	unitOfWork port.UnitOfWork,
	jobPublisher port.JobPublisher,
	requestTokenManager port.EmailChangeRequestTokenManager,
) RequestEmailChangeHandler {
	return &requestEmailChangeHandler{
		userRepo:            userRepo,
		unitOfWork:          unitOfWork,
		jobPublisher:        jobPublisher,
		requestTokenManager: requestTokenManager,
	}
}

func (h *requestEmailChangeHandler) Handle(
	ctx context.Context,
	cmd RequestEmailChange,
) (RequestEmailChangeResult, error) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID(cmd.UserID)

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return RequestEmailChangeResult{}, err
	}
	if user == nil {
		return RequestEmailChangeResult{}, common.ErrUserNotFound
	}

	newEmail, err := vo.NewEmail(cmd.Email)
	if err != nil {
		return RequestEmailChangeResult{}, err
	}
	requestToken, err := h.requestTokenManager.Generate()
	if err != nil {
		return RequestEmailChangeResult{}, err
	}

	if err := user.RequestEmailChange(newEmail, requestToken, common.EmailChangeRequestTTL, now); err != nil {
		return RequestEmailChangeResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.jobPublisher.Publish(ctx, job.NewSendEmail(
			newEmail.String(),
			"emailTemplates/confirm-email-change",
			map[string]string{
				"username": user.Username().String(),
				"token":    user.EmailChangeRequest().Token().String(),
			},
			now,
		)); err != nil {
			return err
		}

		return h.userRepo.Save(ctx, user)
	}); err != nil {
		return RequestEmailChangeResult{}, err
	}

	return RequestEmailChangeResult{}, nil
}
