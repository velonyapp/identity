package port

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrSessionTokenInvalid = errors.New("session token invalid")

type SessionTokenManager interface {
	Generate(ctx context.Context) (vo.SessionToken, error)
	Hash(ctx context.Context, token vo.SessionToken) (vo.SessionTokenHash, error)
	Verify(ctx context.Context, token vo.SessionToken, hash vo.SessionTokenHash) error
}
