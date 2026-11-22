package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type LookupUsers struct {
	UserIDs []string
}

type LookupUsersResult struct {
	Users map[string]common.UserResult
}

func (LookupUsers) resultType() LookupUsersResult {
	return LookupUsersResult{}
}

type LookupUsersHandler Handler[LookupUsers, LookupUsersResult]

type lookupUsersHandler struct {
	userQuery User
}

func NewLookupUsersHandler(
	userQuery User,
) LookupUsersHandler {
	return &lookupUsersHandler{
		userQuery: userQuery,
	}
}

func (h *lookupUsersHandler) Handle(
	ctx context.Context,
	qry LookupUsers,
) (LookupUsersResult, error) {
	userIDs := make([]vo.UserID, len(qry.UserIDs))
	for i, userID := range qry.UserIDs {
		userIDs[i], _ = vo.NewUserID(userID)
	}

	users, err := h.userQuery.BatchFind(ctx, userIDs)
	if err != nil {
		return LookupUsersResult{}, err
	}

	userResults := make(map[string]common.UserResult, len(users))
	for userID, user := range users {
		userResults[userID.String()] = user
	}

	return LookupUsersResult{
		Users: userResults,
	}, nil
}
