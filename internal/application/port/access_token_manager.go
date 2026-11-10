package port

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/vo"
)

type AccessTokenManager interface {
	Generate(ctx context.Context, userID vo.UserID) (string, error)
}
