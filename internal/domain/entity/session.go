package entity

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
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
	ttlSeconds int,
) *Session {
	return &Session{
		ID:         vo.NewSessionIDRandom(),
		UserID:     userID,
		Token:      token,
		ExpireTime: vo.NewTimeNow().AddSeconds(ttlSeconds),
	}
}

func (s *Session) Refresh(token vo.SessionToken, ttlSeconds int) error {
	if s.RevokeTime != nil {
		return ErrSessionRevoked
	}

	now := vo.NewTimeNow()

	if s.ExpireTime.Before(now) {
		return ErrSessionExpired
	}

	s.Token = token
	s.ExpireTime = now.AddSeconds(ttlSeconds)

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
