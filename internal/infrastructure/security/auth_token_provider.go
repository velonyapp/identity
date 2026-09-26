package security

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
)

var _ port.AuthTokenProvider = (*authTokenProvider)(nil)

type authTokenProvider struct {
	c *conf.Security
}

func NewAuthTokenProvider(c *conf.Security) port.AuthTokenProvider {
	return &authTokenProvider{c: c}
}

func (p *authTokenProvider) GenerateAccessToken(subject string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(p.c.GetAccessToken().GetTtl().AsDuration())),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(p.c.GetAccessToken().GetSecret()))
	if err != nil {
		return "", port.ErrGenerateAccessToken
	}

	return signedToken, nil
}

func (p *authTokenProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", port.ErrGenerateRefreshToken
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
