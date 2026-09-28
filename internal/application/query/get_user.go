package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
)

type GetUser struct {
	UserID string
}

type GetUserResult struct {
	User *common.User
}

type GetUserHandler struct {
	userQuery User
	cache     port.Cache
}

func NewGetUserHandler(
	userQuery User,
	cache port.Cache,
) *GetUserHandler {
	return &GetUserHandler{
		userQuery: userQuery,
		cache:     cache,
	}
}

func (h *GetUserHandler) Execute(
	ctx context.Context,
	qry *GetUser,
) (*GetUserResult, error) {
	cacheKey := common.UserCacheKey(qry.UserID)

	cached := &common.User{}

	if found, _ := h.cache.Get(ctx, cacheKey, cached); found {
		return &GetUserResult{User: cached}, nil
	}

	user, err := h.userQuery.GetUser(ctx, qry.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	h.cache.Set(ctx, cacheKey, user, common.UserCacheTTL)

	return &GetUserResult{User: user}, nil
}
