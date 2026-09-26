package command

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type RefreshAuth struct {
	RefreshToken string
}

type RefreshAuthResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshAuthHandler struct {
	c                   *conf.Security
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
	cache               port.Cache
}

func NewRefreshAuthHandler(
	c *conf.Security,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
	cache port.Cache,
) *RefreshAuthHandler {
	return &RefreshAuthHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
		cache:               cache,
	}
}

func (h *RefreshAuthHandler) Execute(
	ctx context.Context,
	cmd *RefreshAuth,
) (*RefreshAuthResult, error) {
	now := time.Now()

	sessionToken, err := vo.NewSessionToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}
	sessionTokenHash, err := h.sessionTokenManager.Hash(sessionToken)
	if err != nil {
		return nil, err
	}

	newSessionToken, err := h.sessionTokenManager.Generate()
	if err != nil {
		return nil, err
	}
	newSessionTokenHash, err := h.sessionTokenManager.Hash(newSessionToken)
	if err != nil {
		return nil, err
	}

	var user *entity.User
	var accessToken, refreshToken string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		session, err := h.sessionRepo.FindByTokenHash(ctx, sessionTokenHash)
		if err != nil {
			return err
		}
		if session == nil {
			return ErrInvalidRefreshToken
		}

		user, err = h.userRepo.FindByID(ctx, session.UserID())
		if err != nil {
			return err
		}
		if user == nil {
			return ErrInvalidRefreshToken
		}

		accessToken, err = h.accessTokenManager.Generate(session.UserID())
		if err != nil {
			return err
		}
		refreshToken = newSessionToken.Value()

		if err := session.Refresh(
			newSessionTokenHash,
			time.Duration(h.c.RefreshToken.Ttl.Seconds)*time.Second,
			now,
		); err != nil {
			return err
		}

		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return nil, err
	}

	h.cache.Set(ctx,
		common.UserResultCacheKey(user.ID()),
		common.NewUserResult(user),
		common.UserResultCacheTTL,
	)

	return &RefreshAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
