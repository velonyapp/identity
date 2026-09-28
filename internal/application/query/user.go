package query

import (
	"context"

	"github.com/velonyapp/identity/internal/application/common"
)

type User interface {
	GetUser(ctx context.Context, userID string) (*common.User, error)
	BatchGetUser(ctx context.Context, userID []string) ([]*common.User, error)
}
