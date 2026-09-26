package vo

import "time"

type EmailChangeRequest struct {
	email      Email
	time       time.Time
	expireTime time.Time
}

func NewEmailChangeRequest(
	email Email,
	time time.Time,
	expireTime time.Time,
) EmailChangeRequest {
	return EmailChangeRequest{
		email:      email,
		time:       time,
		expireTime: expireTime,
	}
}

func (r EmailChangeRequest) Email() Email {
	return r.email
}

func (r EmailChangeRequest) Time() time.Time {
	return r.time
}

func (r EmailChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}
