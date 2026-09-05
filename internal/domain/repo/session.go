package repo

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type Session interface {
	FindByToken(ctx context.Context, username vo.SessionToken) (*entity.Session, error)

	Save(ctx context.Context, user *entity.Session) error
}
