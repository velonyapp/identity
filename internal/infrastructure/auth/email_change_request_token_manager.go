package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.EmailChangeRequestTokenManager = (*emailChangeRequestTokenManager)(nil)

type emailChangeRequestTokenManager struct{}

func NewEmailChangeRequestTokenManager() port.EmailChangeRequestTokenManager {
	return &emailChangeRequestTokenManager{}
}

func (m *emailChangeRequestTokenManager) Generate() (vo.EmailChangeRequestToken, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return vo.EmailChangeRequestToken{}, err
	}

	return vo.NewEmailChangeRequestToken(fmt.Sprintf("%06d", n.Int64()))
}

func (m *emailChangeRequestTokenManager) Verify(provided, expected vo.EmailChangeRequestToken) error {
	if subtle.ConstantTimeCompare(provided.Bytes(), expected.Bytes()) != 1 {
		return port.ErrEmailChangeRequestTokenInvalid
	}

	return nil
}
