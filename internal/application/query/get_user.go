package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type GetUser struct {
	UserID string
}

type GetUserResult struct {
	User common.UserResult
}

func (GetUser) resultType() GetUserResult {
	return GetUserResult{}
}

type GetUserHandler Handler[GetUser, GetUserResult]

type getUserHandler struct {
	userQuery User
}

func NewGetUserHandler(
	userQuery User,
) GetUserHandler {
	return &getUserHandler{
		userQuery: userQuery,
	}
}

func (h *getUserHandler) Handle(
	ctx context.Context,
	qry GetUser,
) (GetUserResult, error) {
	userID, _ := vo.NewUserID(qry.UserID)

	user, err := h.userQuery.Find(ctx, userID)
	if err != nil {
		return GetUserResult{}, err
	}
	if user == nil {
		return GetUserResult{}, common.ErrUserNotFound
	}

	return GetUserResult{
		User: *user,
	}, nil
}
