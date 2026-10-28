package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type User interface {
	Find(ctx context.Context, userID vo.UserID) (*common.UserResult, error)
	BatchFind(ctx context.Context, userIDs []vo.UserID) (map[vo.UserID]common.UserResult, error)
}
