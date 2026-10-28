package port

import "github.com/velonyapp/identity/internal/domain/vo"

type AccessTokenManager interface {
	Generate(userID vo.UserID) (string, error)
}
