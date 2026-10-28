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

type RefreshAuth struct {
	RefreshToken string
}

type RefreshAuthResult struct {
	AccessToken  string
	RefreshToken string
}

func (RefreshAuth) resultType() RefreshAuthResult {
	return RefreshAuthResult{}
}

type RefreshAuthHandler Handler[RefreshAuth, RefreshAuthResult]

type refreshAuthHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewRefreshAuthHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) RefreshAuthHandler {
	return &refreshAuthHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *refreshAuthHandler) Handle(
	ctx context.Context,
	cmd RefreshAuth,
) (RefreshAuthResult, error) {
	now := time.Now().UTC()

	token, _ := vo.NewSessionToken(cmd.RefreshToken)
	tokenHash, err := h.sessionTokenManager.Hash(token)
	if err != nil {
		return RefreshAuthResult{}, err
	}

	session, err := h.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return RefreshAuthResult{}, err
	}
	if session == nil {
		return RefreshAuthResult{}, ErrInvalidRefreshToken
	}

	newSessionToken, err := h.sessionTokenManager.Generate()
	if err != nil {
		return RefreshAuthResult{}, err
	}
	newSessionTokenHash, err := h.sessionTokenManager.Hash(newSessionToken)
	if err != nil {
		return RefreshAuthResult{}, err
	}

	if err := session.Refresh(newSessionTokenHash, h.c.RefreshToken.Ttl.AsDuration(), now); err != nil {
		return RefreshAuthResult{}, err
	}

	newAccessToken, err := h.accessTokenManager.Generate(session.UserID())
	if err != nil {
		return RefreshAuthResult{}, err
	}
	newRefeshToken := newSessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return RefreshAuthResult{}, err
	}

	return RefreshAuthResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefeshToken,
	}, nil
}
