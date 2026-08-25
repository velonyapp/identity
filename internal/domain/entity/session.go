package entity

import (
	"errors"

	"github.com/velony-app/identity/internal/domain/vo"
)

var (
	ErrSessionRevoked = errors.New("session revoked")
	ErrSessionExpired = errors.New("session expired")
)

type Session struct {
	ID         vo.SessionID
	UserID     vo.UserID
	Token      vo.SessionToken
	ExpireTime vo.Time
	RevokeTime *vo.Time
}

func NewSession(
	userID vo.UserID,
	token vo.SessionToken,
) *Session {
	return &Session{
		ID:         vo.NewSessionIDRandom(),
		UserID:     userID,
		Token:      token,
		ExpireTime: vo.NewTimeNow().AddHours(24),
	}
}

func (s *Session) Refresh(token vo.SessionToken) error {
	if s.RevokeTime != nil {
		return ErrSessionRevoked
	}

	now := vo.NewTimeNow()

	if s.ExpireTime.Before(now) {
		return ErrSessionExpired
	}

	s.Token = token
	s.ExpireTime = now.AddHours(24)

	return nil
}

func (s *Session) Revoke() error {
	if s.RevokeTime != nil {
		return ErrSessionRevoked
	}

	now := vo.NewTimeNow()
	s.RevokeTime = &now

	return nil
}
