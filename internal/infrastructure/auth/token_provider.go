package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
)

type TokenProvider struct {
	c *conf.Auth
}

func NewTokenProvider(c *conf.Auth) port.TokenProvider {
	return &TokenProvider{c: c}
}

func (p *TokenProvider) GenerateAccessToken(subject string) (string, error) {
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

func (p *TokenProvider) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", port.ErrGenerateRefreshToken
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
