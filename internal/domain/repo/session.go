package repo

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrSessionVersionConflict = errors.New("session version conflict")
)

type Session interface {
	FindByTokenHash(ctx context.Context, tokenHash vo.SessionTokenHash) (*entity.Session, error)

	Save(ctx context.Context, session *entity.Session) error
}
