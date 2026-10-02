package port

import (
	"context"
	"time"

	"github.com/velonyapp/identity/internal/domain/vo"
)

type AssetService interface {
	CreateAvatar(ctx context.Context, userID vo.UserID) (vo.AvatarID, error)
	PresignAvatar(ctx context.Context, avatarID vo.AvatarID, ttl time.Duration) (string, error)
}
