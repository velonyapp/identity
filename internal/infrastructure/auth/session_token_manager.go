package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.SessionTokenManager = (*sessionTokenManager)(nil)

type sessionTokenManager struct {
	log *slog.Logger
}

func NewSessionTokenManager(
	log *slog.Logger,
) port.SessionTokenManager {
	return &sessionTokenManager{
		log: log,
	}
}

func (m *sessionTokenManager) Generate(ctx context.Context) (vo.SessionToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		m.log.ErrorContext(ctx,
			"failed to generate session token bytes",
			slog.Any("error", err),
		)
		return vo.SessionToken{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	sessionToken, err := vo.NewSessionToken(token)
	if err != nil {
		m.log.ErrorContext(ctx,
			"failed to create session token value object",
			slog.Any("error", err),
		)
		return vo.SessionToken{}, err
	}

	return sessionToken, nil
}

func (m *sessionTokenManager) Hash(ctx context.Context, token vo.SessionToken) (vo.SessionTokenHash, error) {
	hash := sha256.Sum256([]byte(token.String()))

	sessionTokenHash, err := vo.NewSessionTokenHash(hash[:])
	if err != nil {
		m.log.ErrorContext(ctx,
			"failed to create session token hash value object",
			slog.Any("error", err),
		)
		return vo.SessionTokenHash{}, err
	}

	return sessionTokenHash, nil
}

func (m *sessionTokenManager) Verify(ctx context.Context, token vo.SessionToken, hash vo.SessionTokenHash) error {
	computedHash := sha256.Sum256([]byte(token.String()))

	if subtle.ConstantTimeCompare(computedHash[:], hash.Bytes()) != 1 {
		return port.ErrSessionTokenInvalid
	}

	return nil
}
