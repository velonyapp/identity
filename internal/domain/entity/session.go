package entity

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/event"
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

	domainEvents []event.DomainEvent
}

func NewSession(
	userID vo.UserID,
	token vo.SessionToken,
	ttlSeconds int,
) *Session {
	now := vo.NewTimeNow()
	sessionID := vo.NewSessionIDRandom()

	session := &Session{
		ID:         sessionID,
		UserID:     userID,
		Token:      token,
		ExpireTime: vo.NewTimeNow().AddSeconds(ttlSeconds),
	}

	session.recordEvent(
		event.NewSessionCreated(
			sessionID,
			userID,
			now,
		),
	)

	return session
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

	s.recordEvent(
		event.NewSessionRefreshed(
			s.ID,
			now,
		),
	)

	return nil
}

func (s *Session) Revoke() error {
	if s.RevokeTime != nil {
		return ErrSessionRevoked
	}

	now := vo.NewTimeNow()
	s.RevokeTime = &now

	s.recordEvent(
		event.NewSessionRevoked(
			s.ID,
			now,
		),
	)

	return nil
}

func (s *Session) PullEvents() []event.DomainEvent {
	pulled := s.domainEvents
	s.domainEvents = nil
	return pulled
}

func (s *Session) recordEvent(domainEvent event.DomainEvent) {
	s.domainEvents = append(s.domainEvents, domainEvent)
}
