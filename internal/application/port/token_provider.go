package port

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrGenerateAccessToken  = errors.New("failed to generate access token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrGenerateRequestToken = errors.New("failed to generate request token")
)

type TokenProvider interface {
	GenerateAccessToken(subject string) (string, error)
	GenerateRefreshToken() (string, error)
	GenerateRequestToken() (vo.RequestToken, error)
}
