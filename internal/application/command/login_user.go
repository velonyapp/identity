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

type LoginUser struct {
	Identity string
	Password string
}

type LoginUserResult struct {
	AccessToken  string
	RefreshToken string
}

func (LoginUser) resultType() LoginUserResult {
	return LoginUserResult{}
}

type LoginUserHandler Handler[LoginUser, LoginUserResult]

type loginUserHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	passwordManager     port.PasswordManager
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewLoginUserHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	passwordManager port.PasswordManager,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) LoginUserHandler {
	return &loginUserHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		passwordManager:     passwordManager,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *loginUserHandler) Handle(
	ctx context.Context,
	cmd LoginUser,
) (LoginUserResult, error) {
	now := time.Now().UTC()

	username, _ := vo.NewUsername(cmd.Identity)
	email, _ := vo.NewEmail(cmd.Identity)

	user, err := h.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return LoginUserResult{}, err
	}
	if user == nil {
		user, err = h.userRepo.FindByEmail(ctx, email)
		if err != nil {
			return LoginUserResult{}, err
		}
	}
	if user == nil {
		return LoginUserResult{}, ErrInvalidCredentials
	}

	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return LoginUserResult{}, err
	}

	if err := h.passwordManager.Verify(ctx, password, user.LocalAuthStrategy().PasswordHash()); err != nil {
		return LoginUserResult{}, err
	}

	sessionToken, err := h.sessionTokenManager.Generate(ctx)
	if err != nil {
		return LoginUserResult{}, err
	}
	sessionTokenHash, err := h.sessionTokenManager.Hash(ctx, sessionToken)
	if err != nil {
		return LoginUserResult{}, err
	}

	session := entity.NewSession(
		user.ID(),
		sessionTokenHash,
		h.c.RefreshToken.Ttl.AsDuration(),
		now,
	)

	accessToken, err := h.accessTokenManager.Generate(ctx, user.ID())
	if err != nil {
		return LoginUserResult{}, err
	}
	refeshToken := sessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return LoginUserResult{}, err
	}

	return LoginUserResult{
		AccessToken:  accessToken,
		RefreshToken: refeshToken,
	}, nil
}
