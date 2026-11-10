package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.EmailChangeRequestTokenManager = (*emailChangeRequestTokenManager)(nil)

type emailChangeRequestTokenManager struct {
	log *slog.Logger
}

func NewEmailChangeRequestTokenManager(
	log *slog.Logger,
) port.EmailChangeRequestTokenManager {
	return &emailChangeRequestTokenManager{
		log: log,
	}
}

func (m *emailChangeRequestTokenManager) Generate(ctx context.Context) (vo.EmailChangeRequestToken, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		m.log.ErrorContext(ctx,
			"failed to generate email change request token",
			slog.Any("error", err),
		)
		return vo.EmailChangeRequestToken{}, err
	}

	return vo.NewEmailChangeRequestToken(fmt.Sprintf("%06d", n.Int64()))
}

func (m *emailChangeRequestTokenManager) Verify(ctx context.Context, provided, expected vo.EmailChangeRequestToken) error {
	if subtle.ConstantTimeCompare(provided.Bytes(), expected.Bytes()) != 1 {
		return port.ErrEmailChangeRequestTokenInvalid
	}

	return nil
}
