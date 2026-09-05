package api

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/api/v1"
	"github.com/velonyapp/identity/internal/application/command"
	"github.com/velonyapp/identity/internal/application/query"

	kratosjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	userMaximumBatchSize = 100

	userResourcePattern = "users/{user}"
)

type Service struct {
	v1.UnimplementedIdentityServiceServer

	getUserHandler       *query.GetUserHandler
	batchGetUsersHandler *query.BatchGetUsersHandler
	registerAuthHandler  *command.RegisterAuthHandler
	loginAuthHandler     *command.LoginAuthHandler
	refreshAuthHandler   *command.RefreshAuthHandler
	updateUserHandler    *command.UpdateUserHandler
	deleteUserHandler    *command.DeleteUserHandler
}

func NewService(
	getUserHandler *query.GetUserHandler,
	batchGetUsersHandler *query.BatchGetUsersHandler,
	registerAuthHandler *command.RegisterAuthHandler,
	loginAuthHandler *command.LoginAuthHandler,
	refreshAuthHandler *command.RefreshAuthHandler,
	updateUserHandler *command.UpdateUserHandler,
	deleteUserHandler *command.DeleteUserHandler,
) *Service {
	return &Service{
		getUserHandler:       getUserHandler,
		batchGetUsersHandler: batchGetUsersHandler,
		registerAuthHandler:  registerAuthHandler,
		loginAuthHandler:     loginAuthHandler,
		refreshAuthHandler:   refreshAuthHandler,
		updateUserHandler:    updateUserHandler,
		deleteUserHandler:    deleteUserHandler,
	}
}

func (s *Service) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.User, error) {
	var userID string

	if err := resourcename.Sscan(req.GetName(), userResourcePattern, &userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if isSelf(userID) {
		subject, err := subjectFromContext(ctx)
		if err != nil {
			return nil, err
		}

		userID = subject
	}

	result, err := s.getUserHandler.Execute(ctx, &query.GetUser{UserID: userID})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.User{
		Name:      resourcename.Sprint(userResourcePattern, result.User.ID),
		Username:  result.User.Username,
		Email:     result.User.Email,
		FullName:  result.User.FullName,
		AvatarKey: result.User.AvatarKey,
	}, nil
}

func (s *Service) BatchGetUsers(ctx context.Context, req *v1.BatchGetUsersRequest) (*v1.BatchGetUsersResponse, error) {
	names := req.GetNames()

	if len(names) > userMaximumBatchSize {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"names must contain no more than %d entries",
			userMaximumBatchSize,
		)
	}

	userIDs := make([]string, 0, len(names))

	for _, name := range names {
		var userID string

		if err := resourcename.Sscan(name, userResourcePattern, &userID); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		userIDs = append(userIDs, userID)
	}

	result, err := s.batchGetUsersHandler.Execute(ctx, &query.BatchGetUsers{UserIDs: userIDs})
	if err != nil {
		return nil, mapError(err)
	}

	users := make([]*v1.User, len(result.Users))

	for i, user := range result.Users {
		users[i] = &v1.User{
			Name:      resourcename.Sprint(userResourcePattern, user.ID),
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			AvatarKey: user.AvatarKey,
		}
	}

	return &v1.BatchGetUsersResponse{Users: users}, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.User, error) {
	subject, err := subjectFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user := req.GetUser()

	var userID string

	if err := resourcename.Sscan(user.GetName(), userResourcePattern, &userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, status.Error(codes.PermissionDenied, "cannot update another user")
	}

	var username, fullName *string

	paths := req.GetUpdateMask().GetPaths()

	if paths == nil {
		if user.GetUsername() != "" {
			value := user.GetUsername()
			username = &value
		}
		if user.GetFullName() != "" {
			value := user.GetFullName()
			fullName = &value
		}
	} else {
		for _, path := range paths {
			switch path {
			case "*":
				value := user.GetUsername()
				username = &value

				value = user.GetFullName()
				fullName = &value

			case "username":
				value := user.GetUsername()
				username = &value

			case "full_name":
				value := user.GetFullName()
				fullName = &value

			default:
				return nil, status.Errorf(
					codes.InvalidArgument,
					"invalid update mask path %q",
					path,
				)
			}
		}
	}

	result, err := s.updateUserHandler.Execute(ctx, &command.UpdateUser{
		UserID:   userID,
		Username: username,
		FullName: fullName,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.User{
		Name:      resourcename.Sprint(userResourcePattern, result.User.ID),
		Username:  result.User.Username,
		Email:     result.User.Email,
		FullName:  result.User.FullName,
		AvatarKey: result.User.AvatarKey,
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	subject, err := subjectFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var userID string

	if err := resourcename.Sscan(req.GetName(), userResourcePattern, &userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if isSelf(userID) {
		userID = subject
	}

	if userID != subject {
		return nil, status.Error(codes.PermissionDenied, "cannot delete another user")
	}

	if _, err := s.deleteUserHandler.Execute(ctx, &command.DeleteUser{
		UserID: userID,
	}); err != nil {
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) RegisterAuth(ctx context.Context, req *v1.RegisterAuthRequest) (*v1.RegisterAuthResponse, error) {
	result, err := s.registerAuthHandler.Execute(ctx, &command.RegisterAuth{
		FullName: req.GetFullName(),
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.RegisterAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *Service) LoginAuth(ctx context.Context, req *v1.LoginAuthRequest) (*v1.LoginAuthResponse, error) {
	result, err := s.loginAuthHandler.Execute(ctx, &command.LoginAuth{
		Identity: req.GetIdentity(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.LoginAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *Service) RefreshAuth(ctx context.Context, req *v1.RefreshAuthRequest) (*v1.RefreshAuthResponse, error) {
	result, err := s.refreshAuthHandler.Execute(ctx, &command.RefreshAuth{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &v1.RefreshAuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func isSelf(userID string) bool {
	return userID == "me"
}

func subjectFromContext(ctx context.Context) (string, error) {
	claims, ok := kratosjwt.FromContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing authentication")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid subject claim")
	}
	if subject == "" {
		return "", status.Error(codes.Unauthenticated, "missing subject claim")
	}

	return subject, nil
}
