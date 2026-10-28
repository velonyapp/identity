package command

import (
	"context"
	"errors"
	"time"

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

func (LoginAuth) resultType() LoginAuthResult {
	return LoginAuthResult{}
}

type LoginAuthHandler Handler[LoginAuth, LoginAuthResult]

type loginAuthHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	passwordManager     port.PasswordManager
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewLoginAuthHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	passwordManager port.PasswordManager,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) LoginAuthHandler {
	return &loginAuthHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		passwordManager:     passwordManager,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *loginAuthHandler) Handle(
	ctx context.Context,
	cmd LoginAuth,
) (LoginAuthResult, error) {
	now := time.Now().UTC()

	username, _ := vo.NewUsername(cmd.Identity)
	email, _ := vo.NewEmail(cmd.Identity)

	user, err := h.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return LoginAuthResult{}, err
	}
	if user == nil {
		user, err = h.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return LoginAuthResult{}, err
		}
	}
	if user == nil {
		return LoginAuthResult{}, ErrInvalidCredentials
	}

	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return LoginAuthResult{}, err
	}

	if err := h.passwordManager.Verify(password, user.LocalAuthStrategy().PasswordHash()); err != nil {
		return LoginAuthResult{}, err
	}

	sessionToken, err := h.sessionTokenManager.Generate()
	if err != nil {
		return LoginAuthResult{}, err
	}
	sessionTokenHash, err := h.sessionTokenManager.Hash(sessionToken)
	if err != nil {
		return LoginAuthResult{}, err
	}

	session := entity.NewSession(
		user.ID(),
		sessionTokenHash,
		h.c.RefreshToken.Ttl.AsDuration(),
		now,
	)

	accessToken, err := h.accessTokenManager.Generate(user.ID())
	if err != nil {
		return LoginAuthResult{}, err
	}
	refeshToken := sessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return LoginAuthResult{}, err
	}

	return LoginAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refeshToken,
	}, nil
}
