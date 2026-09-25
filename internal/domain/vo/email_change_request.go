package vo

import "time"

type EmailChangeRequest struct {
	value      Email
	time       time.Time
	expireTime time.Time
}

func NewEmailChangeRequest(
	value Email,
	time time.Time,
	expireTime time.Time,
) EmailChangeRequest {
	return EmailChangeRequest{
		value:      value,
		time:       time,
		expireTime: expireTime,
	}
}

func (r EmailChangeRequest) Value() Email {
	return r.value
}

func (r EmailChangeRequest) Time() time.Time {
	return r.time
}

func (r EmailChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}
