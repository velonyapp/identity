package port

import (
	"context"
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrEmailChangeRequestTokenInvalid = errors.New("email change request token invalid")

type EmailChangeRequestTokenManager interface {
	Generate(ctx context.Context) (vo.EmailChangeRequestToken, error)
	Verify(ctx context.Context, provided, expected vo.EmailChangeRequestToken) error
}
