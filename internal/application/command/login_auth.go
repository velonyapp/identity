package command

import (
	"context"
	"errors"

	"github.com/velony-app/identity/internal/application/common"
	"github.com/velony-app/identity/internal/application/port"
	"github.com/velony-app/identity/internal/domain/entity"
	"github.com/velony-app/identity/internal/domain/repo"
	"github.com/velony-app/identity/internal/domain/vo"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginAuth struct {
	Identity string
	Password string
}

type LoginAuthResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginAuthHandler struct {
	userRepo       repo.User
	sessionRepo    repo.Session
	unitOfWork     port.UnitOfWork
	tokenProvider  port.TokenProvider
	passwordHasher port.PasswordHasher
	cache          port.Cache
}

func NewLoginAuthHandler(
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	tokenProvider port.TokenProvider,
	passwordHasher port.PasswordHasher,
	cache port.Cache,
) *LoginAuthHandler {
	return &LoginAuthHandler{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		unitOfWork:     unitOfWork,
		tokenProvider:  tokenProvider,
		passwordHasher: passwordHasher,
		cache:          cache,
	}
}

func (h *LoginAuthHandler) Execute(
	ctx context.Context,
	cmd *LoginAuth,
) (*LoginAuthResult, error) {
	username, _ := vo.NewUsername(cmd.Identity)
	email, _ := vo.NewEmail(cmd.Identity)
	password, _ := vo.NewPassword(cmd.Password)

	var user *entity.User
	var accessToken, refreshToken string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		var err error

		user, err = h.userRepo.FindByUsername(ctx, username)
		if err != nil {
			return err
		}
		if user == nil {
			user, err = h.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return err
			}
		}

		if user == nil {
			return ErrInvalidCredentials
		}
		if user.LocalAuthStrategy == nil {
			return ErrInvalidCredentials
		}
		if err := h.passwordHasher.Verify(password, user.LocalAuthStrategy.PasswordHash); err != nil {
			return ErrInvalidCredentials
		}

		accessToken, err = h.tokenProvider.GenerateAccessToken(user.ID.Value())
		if err != nil {
			return err
		}
		refreshToken, err = h.tokenProvider.GenerateRefreshToken()
		if err != nil {
			return err
		}

		sessionToken, err := vo.NewSessionToken(refreshToken)
		if err != nil {
			return err
		}

		session := entity.NewSession(user.ID, sessionToken)

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

	return &LoginAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
