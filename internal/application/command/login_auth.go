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
	c                   *conf.Security
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
	passwordManager     port.PasswordManager
	cache               port.Cache
}

func NewLoginAuthHandler(
	c *conf.Security,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
	passwordManager port.PasswordManager,
	cache port.Cache,
) *LoginAuthHandler {
	return &LoginAuthHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
		passwordManager:     passwordManager,
		cache:               cache,
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
		if err := h.passwordManager.Verify(password, user.LocalAuthStrategy().PasswordHash()); err != nil {
			return ErrInvalidCredentials
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

		return h.sessionRepo.Save(ctx, session)
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
