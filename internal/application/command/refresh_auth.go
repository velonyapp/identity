package command

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
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
	userRepo      repo.User
	sessionRepo   repo.Session
	unitOfWork    port.UnitOfWork
	tokenProvider port.TokenProvider
	cache         port.Cache
}

func NewRefreshAuthHandler(
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	tokenProvider port.TokenProvider,
	cache port.Cache,
) *RefreshAuthHandler {
	return &RefreshAuthHandler{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		unitOfWork:    unitOfWork,
		tokenProvider: tokenProvider,
		cache:         cache,
	}
}

func (h *RefreshAuthHandler) Execute(
	ctx context.Context,
	cmd *RefreshAuth,
) (*RefreshAuthResult, error) {
	sessionToken, err := vo.NewSessionToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := h.tokenProvider.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newSessionToken, err := vo.NewSessionToken(refreshToken)
	if err != nil {
		return nil, err
	}

	var user *entity.User
	var accessToken string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		session, err := h.sessionRepo.FindByToken(ctx, sessionToken)
		if err != nil {
			return err
		}
		if session == nil {
			return ErrInvalidRefreshToken
		}

		user, err = h.userRepo.FindByID(ctx, session.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrInvalidRefreshToken
		}

		accessToken, err = h.tokenProvider.GenerateAccessToken(session.UserID.Value())
		if err != nil {
			return err
		}

		if err := session.Refresh(newSessionToken); err != nil {
			return err
		}

		if err := h.sessionRepo.Save(ctx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	h.cache.Set(ctx,
		common.UserResultCacheKey(user.ID),
		common.NewUserResult(user),
		common.UserResultCacheTTL,
	)

	return &RefreshAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
