package service

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/identity/gen/api/v1"
	"github.com/velonyapp/identity/internal/application/command"
	"github.com/velonyapp/identity/internal/application/query"

	kjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrInvalidUserResourceName = errors.New("invalid user resource name")
	ErrInvalidUserBatchSize    = errors.New("cannot fetch more than 100 users at once")

	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("unauthorized")
)

const (
	userResourcePattern  = "users/{user}"
	userMaximumBatchSize = 100
)

type API struct {
	v1.UnimplementedIdentityServiceServer

	commandBus *command.Bus
	queryBus   *query.Bus
}

func NewAPI(
	commandBus *command.Bus,
	queryBus *query.Bus,
) *API {
	return &API{
		commandBus: commandBus,
		queryBus:   queryBus,
	}
}

func (s *API) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		subject, err := parseSubject(ctx)
		if err != nil {
			return nil, err
		}

		userID = subject
	}

	result, err := query.Send(ctx, s.queryBus, query.GetUser{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.User{
		Name:       formatUserName(result.User.ID),
		Username:   result.User.Username,
		Email:      result.User.Email,
		AvatarKey:  result.User.AvatarKey,
		CreateTime: timestamppb.New(result.User.CreateTime),
	}, nil
}

func (s *API) BatchGetUsers(ctx context.Context, req *v1.BatchGetUsersRequest) (*v1.BatchGetUsersResponse, error) {
	if len(req.Names) > userMaximumBatchSize {
		return nil, ErrInvalidUserBatchSize
	}

	userIDs := make([]string, len(req.Names))
	for i, name := range req.Names {
		userID, err := parseUserName(name)
		if err != nil {
			return nil, err
		}

		userIDs[i] = userID
	}

	result, err := query.Send(ctx, s.queryBus, query.BatchGetUsers{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	users := make([]*v1.User, len(result.Users))
	for i, user := range result.Users {
		users[i] = &v1.User{
			Name:       formatUserName(user.ID),
			Username:   user.Username,
			Email:      user.Email,
			AvatarKey:  user.AvatarKey,
			CreateTime: timestamppb.New(user.CreateTime),
		}
	}

	return &v1.BatchGetUsersResponse{
		Users: users,
	}, nil
}

func (s *API) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	if _, err := command.Send(ctx, s.commandBus, command.DeleteUser{
		UserID: userID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *API) ChangeUserUsername(ctx context.Context, req *v1.ChangeUserUsernameRequest) (*v1.User, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	result, err := command.Send(ctx, s.commandBus, command.ChangeUsername{
		UserID:   userID,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &v1.User{
		Name:       formatUserName(result.User.ID),
		Username:   result.User.Username,
		Email:      result.User.Email,
		AvatarKey:  result.User.AvatarKey,
		CreateTime: timestamppb.New(result.User.CreateTime),
	}, nil
}

func (s *API) RequestUserEmailChange(ctx context.Context, req *v1.RequestUserEmailChangeRequest) (*v1.RequestUserEmailChangeResponse, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	if _, err := command.Send(ctx, s.commandBus, command.RequestEmailChange{
		UserID: userID,
		Email:  req.Email,
	}); err != nil {
		return nil, err
	}

	return &v1.RequestUserEmailChangeResponse{}, nil
}

func (s *API) ConfirmUserEmailChange(ctx context.Context, req *v1.ConfirmUserEmailChangeRequest) (*v1.User, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	result, err := command.Send(ctx, s.commandBus, command.ConfirmEmailChange{
		UserID: userID,
		Token:  req.Token,
	})
	if err != nil {
		return nil, err
	}

	return &v1.User{
		Name:       formatUserName(result.User.ID),
		Username:   result.User.Username,
		Email:      result.User.Email,
		AvatarKey:  result.User.AvatarKey,
		CreateTime: timestamppb.New(result.User.CreateTime),
	}, nil
}

func (s *API) RemoveUserEmail(ctx context.Context, req *v1.RemoveUserEmailRequest) (*v1.User, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	result, err := command.Send(ctx, s.commandBus, command.RemoveEmail{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.User{
		Name:       formatUserName(result.User.ID),
		Username:   result.User.Username,
		Email:      result.User.Email,
		AvatarKey:  result.User.AvatarKey,
		CreateTime: timestamppb.New(result.User.CreateTime),
	}, nil
}

func (s *API) RequestUserAvatarChange(ctx context.Context, req *v1.RequestUserAvatarChangeRequest) (*v1.RequestUserAvatarChangeResponse, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	result, err := command.Send(ctx, s.commandBus, command.RequestAvatarChange{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RequestUserAvatarChangeResponse{
		UploadUrl: result.UploadURL,
	}, nil
}

func (s *API) RemoveUserAvatar(ctx context.Context, req *v1.RemoveUserAvatarRequest) (*v1.User, error) {
	userID, err := parseUserName(req.Name)
	if err != nil {
		return nil, err
	}

	subject, err := parseSubject(ctx)
	if err != nil {
		return nil, err
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, ErrForbidden
	}

	result, err := command.Send(ctx, s.commandBus, command.RemoveAvatar{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1.User{
		Name:       formatUserName(result.User.ID),
		Username:   result.User.Username,
		Email:      result.User.Email,
		AvatarKey:  result.User.AvatarKey,
		CreateTime: timestamppb.New(result.User.CreateTime),
	}, nil
}

func (s *API) RegisterAuth(ctx context.Context, req *v1.RegisterAuthRequest) (*v1.RegisterAuthResponse, error) {
	result, err := command.Send(ctx, s.commandBus, command.RegisterAuth{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RegisterAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *API) LoginAuth(ctx context.Context, req *v1.LoginAuthRequest) (*v1.LoginAuthResponse, error) {
	result, err := command.Send(ctx, s.commandBus, command.LoginAuth{
		Identity: req.Identity,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.LoginAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *API) RefreshAuth(ctx context.Context, req *v1.RefreshAuthRequest) (*v1.RefreshAuthResponse, error) {
	result, err := command.Send(ctx, s.commandBus, command.RefreshAuth{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RefreshAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func isSelf(userID string) bool {
	return userID == "me"
}

func parseSubject(ctx context.Context) (string, error) {
	claims, ok := kjwt.FromContext(ctx)
	if !ok {
		return "", ErrUnauthenticated
	}

	subject, err := claims.GetSubject()
	if err != nil || subject == "" {
		return "", ErrUnauthenticated
	}

	return subject, nil
}

func parseUserName(name string) (string, error) {
	var userID string
	if err := resourcename.Sscan(name, userResourcePattern, &userID); err != nil {
		return "", ErrInvalidUserResourceName
	}

	return userID, nil
}

func formatUserName(userID string) string {
	return resourcename.Sprint(userResourcePattern, userID)
}
