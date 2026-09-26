package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var _ port.RequestVerifier = (*requestVerifier)(nil)

type requestVerifier struct {
	c *conf.Security
}

func NewRequestVerifier(c *conf.Security) port.RequestVerifier {
	return &requestVerifier{c: c}
}

func (v *requestVerifier) VerifyEmailChange(request vo.EmailChangeRequest, token string) error {
	mac := hmac.New(
		sha256.New,
		[]byte(v.c.GetEmailChangeRequestToken().GetSecret()),
	)

	_, _ = mac.Write([]byte(request.Email().Value()))
	_, _ = mac.Write([]byte(strconv.FormatInt(request.Time().UnixNano(), 10)))
	_, _ = mac.Write([]byte(strconv.FormatInt(request.ExpireTime().UnixNano(), 10)))

	expectedToken := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedToken), []byte(token)) {
		return port.ErrRequestTokenInvalid
	}

	return nil
}

func (v *requestVerifier) VerifyAvatarChange(request vo.AvatarChangeRequest, token string) error {
	mac := hmac.New(
		sha256.New,
		[]byte(v.c.GetAvatarChangeRequestToken().GetSecret()),
	)

	_, _ = mac.Write([]byte(request.AvatarKey().String()))
	_, _ = mac.Write([]byte(strconv.FormatInt(request.Time().UnixNano(), 10)))
	_, _ = mac.Write([]byte(strconv.FormatInt(request.ExpireTime().UnixNano(), 10)))

	expectedToken := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedToken), []byte(token)) {
		return port.ErrRequestTokenInvalid
	}

	return nil
}
