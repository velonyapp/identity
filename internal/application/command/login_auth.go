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
	c              *conf.Security
	userRepo       repo.User
	sessionRepo    repo.Session
	unitOfWork     port.UnitOfWork
	tokenProvider  port.AuthTokenProvider
	passwordHasher port.PasswordHasher
	cache          port.Cache
}

func NewLoginAuthHandler(
	c *conf.Security,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	tokenProvider port.AuthTokenProvider,
	passwordHasher port.PasswordHasher,
	cache port.Cache,
) *LoginAuthHandler {
	return &LoginAuthHandler{
		c:              c,
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
	now := time.Now()

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
		if !user.HasLocalAuthStrategy() {
			return ErrInvalidCredentials
		}
		if err := h.passwordHasher.Verify(password, user.LocalAuthStrategy().PasswordHash()); err != nil {
			return ErrInvalidCredentials
		}

		accessToken, err = h.tokenProvider.GenerateAccessToken(user.ID().Value())
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

		session := entity.NewSession(
			user.ID(),
			sessionToken,
			time.Duration(h.c.RefreshToken.Ttl.Seconds)*time.Second,
			now,
		)

		if err := h.sessionRepo.Save(ctx, session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	h.cache.Set(ctx,
		common.UserResultCacheKey(user.ID()),
		common.NewUserResult(user),
		common.UserResultCacheTTL,
	)

	return &LoginAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
