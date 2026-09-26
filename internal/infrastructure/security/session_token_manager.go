package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.SessionTokenManager = (*sessionTokenManager)(nil)

type sessionTokenManager struct{}

func NewSessionTokenManager() port.SessionTokenManager {
	return &sessionTokenManager{}
}

func (m *sessionTokenManager) Generate() (vo.SessionToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return vo.SessionToken{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	return vo.NewSessionToken(token)
}

func (m *sessionTokenManager) Hash(token vo.SessionToken) (vo.SessionTokenHash, error) {
	hash := sha256.Sum256([]byte(token.Value()))

	return vo.NewSessionTokenHash(hash[:])
}

func (m *sessionTokenManager) Verify(token vo.SessionToken, hash vo.SessionTokenHash) error {
	computedHash := sha256.Sum256([]byte(token.Value()))

	if subtle.ConstantTimeCompare(computedHash[:], hash.Value()) != 1 {
		return port.ErrSessionTokenInvalid
	}

	return nil
}
