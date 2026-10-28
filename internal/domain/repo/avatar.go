package repo

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type Avatar interface {
	FindByID(ctx context.Context, avatarID vo.AvatarID) (*entity.Avatar, error)

	Create(ctx context.Context, avatar *entity.Avatar) error
	Delete(ctx context.Context, avatarID vo.AvatarID) error
}
