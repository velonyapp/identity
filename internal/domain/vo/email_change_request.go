package vo

import "time"

type EmailChangeRequest struct {
	email      Email
	token      EmailChangeRequestToken
	expireTime time.Time
}

func NewEmailChangeRequest(
	email Email,
	token EmailChangeRequestToken,
	expireTime time.Time,
) EmailChangeRequest {
	return EmailChangeRequest{
		email:      email,
		token:      token,
		expireTime: expireTime,
	}
}

func (r EmailChangeRequest) Email() Email {
	return r.email
}

func (r EmailChangeRequest) Token() EmailChangeRequestToken {
	return r.token
}

func (r EmailChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}

func (r EmailChangeRequest) Equal(other EmailChangeRequest) bool {
	return r.email.Equal(other.email) && r.expireTime.Equal(other.expireTime)
}
