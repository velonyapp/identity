package entity

import (
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/domain/event"
	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrSessionRevoked = errors.New("session revoked")
	ErrSessionExpired = errors.New("session expired")
)

type Session struct {
	id         vo.SessionID
	userID     vo.UserID
	tokenHash  vo.SessionTokenHash
	expireTime time.Time
	revokeTime *time.Time

	domainEvents []event.DomainEvent
}

func NewSession(
	userID vo.UserID,
	tokenHash vo.SessionTokenHash,
	ttl time.Duration,
	now time.Time,
) *Session {
	sessionID := vo.NewSessionIDRandom()

	session := &Session{
		id:         sessionID,
		userID:     userID,
		tokenHash:  tokenHash,
		expireTime: now.Add(ttl),
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

func ReconstituteSession(
	id vo.SessionID,
	userID vo.UserID,
	tokenHash vo.SessionTokenHash,
	expireTime time.Time,
	revokeTime *time.Time,
) *Session {
	session := &Session{
		id:         id,
		userID:     userID,
		tokenHash:  tokenHash,
		expireTime: expireTime,
	}

	if revokeTime != nil {
		value := *revokeTime
		session.revokeTime = &value
	}

	return session
}

func (s *Session) ID() vo.SessionID {
	return s.id
}

func (s *Session) UserID() vo.UserID {
	return s.userID
}

func (s *Session) TokenHash() vo.SessionTokenHash {
	return s.tokenHash
}

func (s *Session) ExpireTime() time.Time {
	return s.expireTime
}

func (s *Session) RevokeTime() *time.Time {
	if s.revokeTime == nil {
		return nil
	}

	value := *s.revokeTime
	return &value
}

func (s *Session) IsRevoked() bool {
	return s.revokeTime != nil
}

func (s *Session) IsExpired(now time.Time) bool {
	return !now.Before(s.expireTime)
}

func (s *Session) Refresh(
	tokenHash vo.SessionTokenHash,
	ttl time.Duration,
	now time.Time,
) error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}
	if s.IsExpired(now) {
		return ErrSessionExpired
	}

	s.tokenHash = tokenHash
	s.expireTime = now.Add(ttl)

	s.recordEvent(
		event.NewSessionRefreshed(
			s.id,
			now,
		),
	)

	return nil
}

func (s *Session) Revoke(now time.Time) error {
	if s.IsRevoked() {
		return ErrSessionRevoked
	}

	s.revokeTime = &now

	s.recordEvent(
		event.NewSessionRevoked(
			s.id,
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
