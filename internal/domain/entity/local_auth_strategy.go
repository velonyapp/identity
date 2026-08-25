package entity

import "github.com/velony-app/identity/internal/domain/vo"

type LocalAuthStrategy struct {
	UserID       vo.UserID
	PasswordHash vo.PasswordHash
}

func NewLocalAuthStrategy(
	userID vo.UserID,
	passwordHash vo.PasswordHash,
) *LocalAuthStrategy {
	return &LocalAuthStrategy{
		UserID:       userID,
		PasswordHash: passwordHash,
	}
}

func (u *LocalAuthStrategy) ChangePassword(passwordHash vo.PasswordHash) error {
	u.PasswordHash = passwordHash

	return nil
}
