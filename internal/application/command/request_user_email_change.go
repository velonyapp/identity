package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RequestUserEmailChange struct {
	UserID   string
	NewEmail string
}

type RequestUserEmailChangeResult struct{}

type RequestUserEmailChangeHandler struct {
	requestRepo       repo.EmailChangeRequest
	unitOfWork        port.UnitOfWork
	tokenProvider     port.TokenProvider
	emailAvailability *service.EmailAvailability
}

func NewRequestUserEmailChangeHandler(
	requestRepo repo.EmailChangeRequest,
	unitOfWork port.UnitOfWork,
	tokenProvider port.TokenProvider,
	emailAvailability *service.EmailAvailability,
) *RequestUserEmailChangeHandler {
	return &RequestUserEmailChangeHandler{
		requestRepo:       requestRepo,
		unitOfWork:        unitOfWork,
		tokenProvider:     tokenProvider,
		emailAvailability: emailAvailability,
	}
}

func (h *RequestUserEmailChangeHandler) Execute(
	ctx context.Context,
	cmd *RequestUserEmailChange,
) (*RequestUserEmailChangeResult, error) {
	userID := vo.NewUserID(cmd.UserID)
	newEmail, err := vo.NewEmail(cmd.NewEmail)
	if err != nil {
		return nil, err
	}

	token, err := h.tokenProvider.GenerateRequestToken()
	if err != nil {
		return nil, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.emailAvailability.EnsureAvailable(ctx, newEmail); err != nil {
			return err
		}

		request := entity.NewEmailChangeRequest(
			token,
			userID,
			newEmail,
			3600, // TODO: Move this to configuration.
		)

		return h.requestRepo.Save(ctx, request)
	}); err != nil {
		return nil, err
	}

	return &RequestUserEmailChangeResult{}, nil
}
