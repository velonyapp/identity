package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/service"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RegisterAuth struct {
	FullName string
	Username string
	Password string
}

type RegisterAuthResult struct {
	AccessToken  string
	RefreshToken string
}

type RegisterAuthHandler struct {
	c                    *conf.Auth
	userRepo             repo.User
	sessionRepo          repo.Session
	unitOfWork           port.UnitOfWork
	tokenProvider        port.TokenProvider
	passwordHasher       port.PasswordHasher
	cache                port.Cache
	usernameAvailability *service.UsernameAvailability
}

func NewRegisterAuthHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	tokenProvider port.TokenProvider,
	passwordHasher port.PasswordHasher,
	cache port.Cache,
	usernameAvailability *service.UsernameAvailability,
) *RegisterAuthHandler {
	return &RegisterAuthHandler{
		c:                    c,
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		unitOfWork:           unitOfWork,
		tokenProvider:        tokenProvider,
		passwordHasher:       passwordHasher,
		cache:                cache,
		usernameAvailability: usernameAvailability,
	}
}

func (h *RegisterAuthHandler) Execute(
	ctx context.Context,
	cmd *RegisterAuth,
) (*RegisterAuthResult, error) {
	fullName, err := vo.NewFullName(cmd.FullName)
	if err != nil {
		return nil, err
	}
	username, err := vo.NewUsername(cmd.Username)
	if err != nil {
		return nil, err
	}
	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := h.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	var user *entity.User
	var accessToken, refreshToken string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.usernameAvailability.EnsureAvailable(ctx, username); err != nil {
			return err
		}

		user = entity.NewUser(username, nil, fullName, nil)

		localAuthStrategy := entity.NewLocalAuthStrategy(user.ID, passwordHash)

		if err := user.ReplaceLocalAuthStrategy(localAuthStrategy); err != nil {
			return err
		}

		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
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

		session := entity.NewSession(user.ID, sessionToken, int(h.c.RefreshToken.Ttl.Seconds))

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

	return &RegisterAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
