package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/application/port"
)

type BatchGetUsers struct {
	UserIDs []string
}

type BatchGetUsersResult struct {
	Users []*common.User
}

type BatchGetUsersHandler struct {
	userQuery User
	cache     port.Cache
}

func NewBatchGetUsersHandler(
	userQuery User,
	cache port.Cache,
) *BatchGetUsersHandler {
	return &BatchGetUsersHandler{
		userQuery: userQuery,
		cache:     cache,
	}
}

func (h *BatchGetUsersHandler) Execute(
	ctx context.Context,
	qry *BatchGetUsers,
) (*BatchGetUsersResult, error) {
	userIDs := make([]string, 0, len(qry.UserIDs))
	cacheKeys := make([]string, 0, len(qry.UserIDs))
	cached := make(map[string]any, len(qry.UserIDs))
	cachedResults := make(map[string]*common.User, len(qry.UserIDs))

	for _, userID := range qry.UserIDs {
		cacheKey := common.UserCacheKey(userID)
		user := &common.User{}

		userIDs = append(userIDs, userID)
		cacheKeys = append(cacheKeys, cacheKey)

		cached[cacheKey] = user
		cachedResults[cacheKey] = user
	}

	found, _ := h.cache.GetMany(ctx, cacheKeys, cached)

	results := make([]*common.User, len(userIDs))
	missingUserIDs := make([]string, 0)

	for i, userID := range userIDs {
		cacheKey := cacheKeys[i]

		if found[cacheKey] {
			results[i] = cachedResults[cacheKey]
			continue
		}

		missingUserIDs = append(missingUserIDs, userID)
	}

	if len(missingUserIDs) == 0 {
		return &BatchGetUsersResult{Users: results}, nil
	}

	users, err := h.userQuery.BatchGetUser(ctx, missingUserIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[string]*common.User, len(users))
	for _, user := range users {
		usersByID[user.ID] = user
	}

	cacheItems := make(map[string]any, len(users))

	for i, userID := range userIDs {
		if results[i] != nil {
			continue
		}

		user, ok := usersByID[userID]
		if !ok {
			return nil, common.ErrUserNotFound
		}

		results[i] = user
		cacheItems[cacheKeys[i]] = user
	}

	h.cache.SetMany(ctx, cacheItems, common.UserCacheTTL)

	return &BatchGetUsersResult{Users: results}, nil
}
