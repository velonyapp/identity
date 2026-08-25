package query

import (
	"context"

	"github.com/velony-app/identity/internal/application/common"
	"github.com/velony-app/identity/internal/application/port"
	"github.com/velony-app/identity/internal/domain/entity"
	"github.com/velony-app/identity/internal/domain/repo"
	"github.com/velony-app/identity/internal/domain/vo"
)

type BatchGetUsers struct {
	UserIDs []string
}

type BatchGetUsersResult struct {
	Users []*common.UserResult
}

type BatchGetUsersHandler struct {
	userRepo repo.User
	cache    port.Cache
}

func NewBatchGetUsersHandler(
	userRepo repo.User,
	cache port.Cache,
) *BatchGetUsersHandler {
	return &BatchGetUsersHandler{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (h *BatchGetUsersHandler) Execute(
	ctx context.Context,
	qry *BatchGetUsers,
) (*BatchGetUsersResult, error) {
	userIDs := make([]vo.UserID, 0, len(qry.UserIDs))
	cacheKeys := make([]string, 0, len(qry.UserIDs))
	cached := make(map[string]any, len(qry.UserIDs))
	cachedResults := make(map[string]*common.UserResult, len(qry.UserIDs))

	for _, userID := range qry.UserIDs {
		id := vo.NewUserID(userID)
		cacheKey := common.UserResultCacheKey(id)
		userResult := &common.UserResult{}

		userIDs = append(userIDs, id)
		cacheKeys = append(cacheKeys, cacheKey)

		cached[cacheKey] = userResult
		cachedResults[cacheKey] = userResult
	}

	found, _ := h.cache.GetMany(ctx, cacheKeys, cached)

	results := make([]*common.UserResult, len(userIDs))
	missingUserIDs := make([]vo.UserID, 0)

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

	users, err := h.userRepo.FindByIDs(ctx, missingUserIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[vo.UserID]*entity.User, len(users))
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

		userResult := common.NewUserResult(user)

		results[i] = userResult
		cacheItems[cacheKeys[i]] = userResult
	}

	h.cache.SetMany(ctx, cacheItems, common.UserResultCacheTTL)

	return &BatchGetUsersResult{Users: results}, nil
}
