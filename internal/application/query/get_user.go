package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/repo"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type GetUser struct {
	UserID string
}

type GetUserResult struct {
	User *common.UserResult
}

type GetUserHandler struct {
	userRepo repo.User
	cache    port.Cache
}

func NewGetUserHandler(
	userRepo repo.User,
	cache port.Cache,
) *GetUserHandler {
	return &GetUserHandler{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (h *GetUserHandler) Execute(
	ctx context.Context,
	qry *GetUser,
) (*GetUserResult, error) {
	userID := vo.NewUserID(qry.UserID)

	cacheKey := common.UserResultCacheKey(userID)

	cached := &common.UserResult{}

	if found, _ := h.cache.Get(ctx, cacheKey, cached); found {
		return &GetUserResult{User: cached}, nil
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	userResult := common.NewUserResult(user)

	h.cache.Set(ctx, cacheKey, userResult, common.UserResultCacheTTL)

	return &GetUserResult{User: common.NewUserResult(user)}, nil
}
