package port

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var ErrEmailChangeRequestTokenInvalid = errors.New("email change request token invalid")

type EmailChangeRequestTokenManager interface {
	Generate() (vo.EmailChangeRequestToken, error)
	Verify(provided, expected vo.EmailChangeRequestToken) error
}
