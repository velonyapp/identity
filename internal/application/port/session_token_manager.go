package port

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrSessionTokenInvalid = errors.New("session token invalid")

type SessionTokenManager interface {
	Generate() (vo.SessionToken, error)
	Hash(token vo.SessionToken) (vo.SessionTokenHash, error)
	Verify(token vo.SessionToken, hash vo.SessionTokenHash) error
}
