package auth

import (
	"time"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/vo"

	"github.com/golang-jwt/jwt/v5"
)

var _ port.AccessTokenManager = (*accessTokenManager)(nil)

type accessTokenManager struct {
	c *conf.Auth
}

func NewAccessTokenManager(c *conf.Auth) port.AccessTokenManager {
	return &accessTokenManager{c: c}
}

func (p *accessTokenManager) Generate(userID vo.UserID) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(p.c.AccessToken.Ttl.AsDuration())),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(p.c.AccessToken.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
