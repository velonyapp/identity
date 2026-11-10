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

type RegisterUser struct {
	Username string
	Password string
}

type RegisterUserResult struct {
	AccessToken  string
	RefreshToken string
}

func (RegisterUser) resultType() RegisterUserResult {
	return RegisterUserResult{}
}

type RegisterUserHandler Handler[RegisterUser, RegisterUserResult]

type registerUserHandler struct {
	c                   *conf.Auth
	userRepo            repo.User
	sessionRepo         repo.Session
	unitOfWork          port.UnitOfWork
	passwordManager     port.PasswordManager
	accessTokenManager  port.AccessTokenManager
	sessionTokenManager port.SessionTokenManager
}

func NewRegisterUserHandler(
	c *conf.Auth,
	userRepo repo.User,
	sessionRepo repo.Session,
	unitOfWork port.UnitOfWork,
	passwordManager port.PasswordManager,
	accessTokenManager port.AccessTokenManager,
	sessionTokenManager port.SessionTokenManager,
) RegisterUserHandler {
	return &registerUserHandler{
		c:                   c,
		userRepo:            userRepo,
		sessionRepo:         sessionRepo,
		unitOfWork:          unitOfWork,
		passwordManager:     passwordManager,
		accessTokenManager:  accessTokenManager,
		sessionTokenManager: sessionTokenManager,
	}
}

func (h *registerUserHandler) Handle(
	ctx context.Context,
	cmd RegisterUser,
) (RegisterUserResult, error) {
	now := time.Now().UTC()

	username, err := vo.NewUsername(cmd.Username)
	if err != nil {
		return RegisterUserResult{}, err
	}
	password, err := vo.NewPassword(cmd.Password)
	if err != nil {
		return RegisterUserResult{}, err
	}
	passwordHash, err := h.passwordManager.Hash(ctx, password)
	if err != nil {
		return RegisterUserResult{}, err
	}

	user := entity.NewUser(
		username,
		now,
	)
	localAuthStrategy := vo.NewLocalAuthStrategy(passwordHash)
	user.ChangeLocalAuthStrategy(localAuthStrategy)

	sessionToken, err := h.sessionTokenManager.Generate(ctx)
	if err != nil {
		return RegisterUserResult{}, err
	}
	sessionTokenHash, err := h.sessionTokenManager.Hash(ctx, sessionToken)
	if err != nil {
		return RegisterUserResult{}, err
	}

	session := entity.NewSession(
		user.ID(),
		sessionTokenHash,
		h.c.RefreshToken.Ttl.AsDuration(),
		now,
	)

	accessToken, err := h.accessTokenManager.Generate(ctx, user.ID())
	if err != nil {
		return RegisterUserResult{}, err
	}
	refeshToken := sessionToken.String()

	if err := h.unitOfWork.Do(ctx, func(ctx context.Context) error {
		if err := h.userRepo.Save(ctx, user); err != nil {
			return err
		}

		return h.sessionRepo.Save(ctx, session)
	}); err != nil {
		return RegisterUserResult{}, err
	}

	return RegisterUserResult{
		AccessToken:  accessToken,
		RefreshToken: refeshToken,
	}, nil
}
