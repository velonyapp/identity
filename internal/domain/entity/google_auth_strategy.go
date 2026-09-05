package entity

import "github.com/velonyapp/identity/internal/domain/vo"

type GoogleAuthStrategy struct {
	UserID vo.UserID
	Sub    vo.GoogleSub
}

func NewGoogleAuthStrategy(
	userID vo.UserID,
	sub vo.GoogleSub,
) *GoogleAuthStrategy {
	return &GoogleAuthStrategy{
		UserID: userID,
		Sub:    sub,
	}
}
