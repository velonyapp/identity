package entity

import (
	"errors"
	"time"

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
	createTime time.Time
	expireTime time.Time
	revokeTime *time.Time

	version uint32
	isNew   bool
}

func NewSession(
	userID vo.UserID,
	tokenHash vo.SessionTokenHash,
	ttl time.Duration,
	now time.Time,
) *Session {
	return &Session{
		id:         vo.GenerateSessionID(),
		userID:     userID,
		tokenHash:  tokenHash,
		createTime: now,
		expireTime: now.Add(ttl),

		version: 0,
		isNew:   true,
	}
}

func ReconstituteSession(
	id vo.SessionID,
	userID vo.UserID,
	tokenHash vo.SessionTokenHash,
	createTime time.Time,
	expireTime time.Time,
	revokeTime *time.Time,
	version uint32,
) *Session {
	s := &Session{
		id:         id,
		userID:     userID,
		tokenHash:  tokenHash,
		createTime: createTime,
		expireTime: expireTime,
		version:    version,
	}

	if revokeTime != nil {
		value := *revokeTime
		s.revokeTime = &value
	}

	return s
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

func (s *Session) CreateTime() time.Time {
	return s.createTime
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

	return nil
}

func (s *Session) Revoke(now time.Time) error {
	if s.IsRevoked() {
		return nil
	}

	s.revokeTime = &now

	return nil
}

func (s *Session) Version() uint32 {
	return s.version
}

func (s *Session) IsNew() bool {
	return s.isNew
}

func (s *Session) MarkPersisted() {
	if s.isNew{
		s.isNew = false

		return
	}
	
	s.version++
}
