package port

import (
	"errors"
)

var (
	ErrGenerateAccessToken  = errors.New("failed to generate access token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token")
)

type AuthTokenProvider interface {
	GenerateAccessToken(subject string) (string, error)
	GenerateRefreshToken() (string, error)
}
