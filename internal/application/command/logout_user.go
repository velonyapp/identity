package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type LogoutUser struct {
	SessionID string
}

type LogoutUserResult struct{}

func (LogoutUser) resultType() LogoutUserResult {
	return LogoutUserResult{}
}

type LogoutUserHandler Handler[LogoutUser, LogoutUserResult]

type logoutUserHandler struct {
	c           *conf.Auth
	sessionRepo repo.Session
	unitOfWork  port.UnitOfWork
}

func NewLogoutUserHandler(
	c *conf.Auth,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
) LogoutUserHandler {
	return &logoutUserHandler{
		c:           c,
		sessionRepo: sessionRepo,
		unitOfWork:  unitOfWork,
	}
}

func (h *logoutUserHandler) Handle(
	ctx context.Context,
	cmd LogoutUser,
) (LogoutUserResult, error) {
	now := time.Now().UTC()

	sessionID, err := vo.NewSessionID(cmd.SessionID)

	session, err := h.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return LogoutUserResult{}, err
	}
	if session == nil {
		return LogoutUserResult{}, nil
	}

	if err := session.Revoke(now); err != nil {
		return LogoutUserResult{}, err
	}

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return LogoutUserResult{}, err
	}

	return LogoutUserResult{}, nil
}
