package port

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrSessionTokenInvalid = errors.New("session token invalid")

type SessionTokenManager interface {
	Generate() (vo.SessionToken, error)
	Hash(sessionToken vo.SessionToken) (vo.SessionTokenHash, error)
	Verify(sessionToken vo.SessionToken, hash vo.SessionTokenHash) error
}
