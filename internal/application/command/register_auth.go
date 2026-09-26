package command

import (
	"context"
	"time"

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
	c                    *conf.Security
	userRepo             repo.User
	sessionRepo          repo.Session
	unitOfWork           port.UnitOfWork
	accessTokenManager   port.AccessTokenManager
	sessionTokenManager  port.SessionTokenManager
	passwordManager      port.PasswordManager
	cache                port.Cache
	usernameAvailability *service.UsernameAvailability
}

func NewRegisterAuthHandler(
	c *conf.Security,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
	passwordManager port.PasswordManager,
	cache port.Cache,
	usernameAvailability *service.UsernameAvailability,
) *RegisterAuthHandler {
	return &RegisterAuthHandler{
		c:                    c,
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		unitOfWork:           unitOfWork,
		accessTokenManager:   accessTokenManager,
		sessionTokenManager:  sessionTokenManager,
		passwordManager:      passwordManager,
		cache:                cache,
		usernameAvailability: usernameAvailability,
	}
}

func (h *RegisterAuthHandler) Execute(
	ctx context.Context,
	cmd *RegisterAuth,
) (*RegisterAuthResult, error) {
	now := time.Now()

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

	passwordHash, err := h.passwordManager.Hash(password)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(username, fullName, now)

	var accessToken, refreshToken string

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.usernameAvailability.EnsureAvailable(ctx, username); err != nil {
			return err
		}

		localAuthStrategy := vo.NewLocalAuthStrategy(passwordHash)

		if err := user.ChangeLocalAuthStrategy(localAuthStrategy); err != nil {
			return err
		}

		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		sessionToken, err := h.sessionTokenManager.Generate()
		if err != nil {
			return err
		}
		sessionTokenHash, err := h.sessionTokenManager.Hash(sessionToken)
		if err != nil {
			return err
		}

		accessToken, err = h.accessTokenManager.Generate(user.ID())
		if err != nil {
			return err
		}
		refreshToken = sessionToken.Value()

		session := entity.NewSession(
			user.ID(),
			sessionTokenHash,
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

	return &RegisterAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
