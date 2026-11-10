package port

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/vo"
)

type AccessTokenManager interface {
	Generate(ctx context.Context, userID vo.UserID, sessionID vo.SessionID) (string, error)
}
