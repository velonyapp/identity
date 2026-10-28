package command

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type RegisterAuth struct {
	Username string
	Password string
}

type RegisterAuthResult struct {
	AccessToken  string
	RefreshToken string
}

func (RegisterAuth) resultType() RegisterAuthResult {
	return RegisterAuthResult{}
}

type RegisterAuthHandler Handler[RegisterAuth, RegisterAuthResult]

type registerAuthHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	passwordManager     port.PasswordManager
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewRegisterAuthHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	passwordManager port.PasswordManager,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) RegisterAuthHandler {
	return &registerAuthHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		passwordManager:     passwordManager,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *registerAuthHandler) Handle(
	ctx context.Context,
	cmd RegisterAuth,
) (RegisterAuthResult, error) {
	now := time.Now().UTC()

	username, err := vo.NewUsername(cmd.Username)
	if err != nil {
		return RegisterAuthResult{}, err
	}
	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return RegisterAuthResult{}, err
	}
	passwordHash, err := h.passwordManager.Hash(password)
	if err != nil {
		return RegisterAuthResult{}, err
	}

	user := entity.NewUser(
		username,
		now,
	)
	localAuthStrategy := vo.NewLocalAuthStrategy(passwordHash)
	user.ChangeLocalAuthStrategy(localAuthStrategy)

	sessionToken, err := h.sessionTokenManager.Generate()
	if err != nil {
		return RegisterAuthResult{}, err
	}
	sessionTokenHash, err := h.sessionTokenManager.Hash(sessionToken)
	if err != nil {
		return RegisterAuthResult{}, err
	}

	session := entity.NewSession(
		user.ID(),
		sessionTokenHash,
		h.c.RefreshToken.Ttl.AsDuration(),
		now,
	)

	accessToken, err := h.accessTokenManager.Generate(user.ID())
	if err != nil {
		return RegisterAuthResult{}, err
	}
	refeshToken := sessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return RegisterAuthResult{}, err
	}

	return RegisterAuthResult{
		AccessToken:  accessToken,
		RefreshToken: refeshToken,
	}, nil
}
