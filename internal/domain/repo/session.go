package repo

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type Session interface {
	FindByTokenHash(ctx context.Context, tokenHash vo.SessionTokenHash) (*entity.Session, error)

	Save(ctx context.Context, user *entity.Session) error
}
