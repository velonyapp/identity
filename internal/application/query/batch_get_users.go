package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type BatchGetUsers struct {
	UserIDs []string
}

type BatchGetUsersResult struct {
	Users []common.UserResult
}

func (BatchGetUsers) resultType() BatchGetUsersResult {
	return BatchGetUsersResult{}
}

type BatchGetUsersHandler Handler[BatchGetUsers, BatchGetUsersResult]

type batchGetUsersHandler struct {
	userQuery User
}

func NewBatchGetUsersHandler(
	userQuery User,
) BatchGetUsersHandler {
	return &batchGetUsersHandler{
		userQuery: userQuery,
	}
}

func (h *batchGetUsersHandler) Handle(
	ctx context.Context,
	qry BatchGetUsers,
) (BatchGetUsersResult, error) {
	userIDs := make([]vo.UserID, len(qry.UserIDs))
	for i, userID := range qry.UserIDs {
		userIDs[i], _ = vo.NewUserID(userID)
	}

	users, err := h.userQuery.BatchFind(ctx, userIDs)
	if err != nil {
		return BatchGetUsersResult{}, err
	}

	userResults := make([]common.UserResult, len(userIDs))
	for i, userID := range userIDs {
		user, ok := users[userID]
		if !ok {
			return BatchGetUsersResult{}, common.ErrUserNotFound
		}

		userResults[i] = user
	}

	return BatchGetUsersResult{
		Users: userResults,
	}, nil
}
