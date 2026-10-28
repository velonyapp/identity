package port

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type AssetService interface {
	CreateAvatar(ctx context.Context, userID vo.UserID) (*entity.Avatar, string, error)
}
