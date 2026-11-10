package command

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type RefreshUserSession struct {
	RefreshToken string
}

type RefreshUserSessionResult struct {
	AccessToken  string
	RefreshToken string
}

func (RefreshUserSession) resultType() RefreshUserSessionResult {
	return RefreshUserSessionResult{}
}

type RefreshUserSessionHandler Handler[RefreshUserSession, RefreshUserSessionResult]

type refreshUserSessionHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewRefreshUserSessionHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) RefreshUserSessionHandler {
	return &refreshUserSessionHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *refreshUserSessionHandler) Handle(
	ctx context.Context,
	cmd RefreshUserSession,
) (RefreshUserSessionResult, error) {
	now := time.Now().UTC()

	token, _ := vo.NewSessionToken(cmd.RefreshToken)
	tokenHash, err := h.sessionTokenManager.Hash(ctx, token)
	if err != nil {
		return RefreshUserSessionResult{}, err
	}

	session, err := h.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return RefreshUserSessionResult{}, err
	}
	if session == nil {
		return RefreshUserSessionResult{}, ErrInvalidRefreshToken
	}

	newSessionToken, err := h.sessionTokenManager.Generate(ctx)
	if err != nil {
		return RefreshUserSessionResult{}, err
	}
	newSessionTokenHash, err := h.sessionTokenManager.Hash(ctx, newSessionToken)
	if err != nil {
		return RefreshUserSessionResult{}, err
	}

	if err := session.Refresh(newSessionTokenHash, h.c.RefreshToken.Ttl.AsDuration(), now); err != nil {
		return RefreshUserSessionResult{}, err
	}

	newAccessToken, err := h.accessTokenManager.Generate(ctx, session.UserID(), session.ID())
	if err != nil {
		return RefreshUserSessionResult{}, err
	}
	newRefeshToken := newSessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return RefreshUserSessionResult{}, err
	}

	return RefreshUserSessionResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefeshToken,
	}, nil
}
