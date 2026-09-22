package entity

import (
	"errors"

	"github.com/velonyapp/identity/internal/domain/vo"
)

var (
	ErrEmailChangeRequestConfirmed = errors.New("email change request confirmed")
	ErrEmailChangeRequestExpired   = errors.New("email change request expired")
)

type EmailChangeRequest struct {
	ID           vo.EmailChangeRequestID
	PendingToken *vo.RequestToken
	UserID       vo.UserID
	NewEmail     vo.Email
	CreateTime   vo.Time
	ExpireTime   vo.Time
	ConfirmTime  *vo.Time
}

func NewEmailChangeRequest(
	token vo.RequestToken,
	userID vo.UserID,
	newEmail vo.Email,
	ttlSeconds int,
) *EmailChangeRequest {
	now := vo.NewTimeNow()
	requestID := vo.NewEmailChangeRequestIDRandom()

	return &EmailChangeRequest{
		ID:           requestID,
		PendingToken: &token,
		UserID:       userID,
		NewEmail:     newEmail,
		CreateTime:   now,
		ExpireTime:   now.AddSeconds(ttlSeconds),
	}
}

func (r *EmailChangeRequest) Confirm() error {
	if r.ConfirmTime != nil {
		return ErrEmailChangeRequestConfirmed
	}

	now := vo.NewTimeNow()

	if r.ExpireTime.Before(now) {
		return ErrEmailChangeRequestExpired
	}

	r.ConfirmTime = &now

	return nil
}
