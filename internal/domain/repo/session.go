package repo

import (
	"context"

	"github.com/velony-app/identity/internal/domain/entity"
	"github.com/velony-app/identity/internal/domain/vo"
)

type Session interface {
	FindByToken(ctx context.Context, username vo.SessionToken) (*entity.Session, error)

	Save(ctx context.Context, user *entity.Session) error
}
