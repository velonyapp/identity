package port

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrRequestTokenInvalid = errors.New("request token invalid")
)

type RequestVerifier interface {
	VerifyEmailChange(request vo.EmailChangeRequest, token string) error
	VerifyAvatarChange(request vo.AvatarChangeRequest, token string) error
}
