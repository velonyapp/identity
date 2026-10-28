package vo

import "time"

type EmailChangeRequest struct {
	email      Email
	token      EmailChangeRequestToken
	createTime time.Time
	expireTime time.Time
}

func NewEmailChangeRequest(
	email Email,
	token EmailChangeRequestToken,
	createTime time.Time,
	expireTime time.Time,
) EmailChangeRequest {
	return EmailChangeRequest{
		email:      email,
		token:      token,
		createTime: createTime,
		expireTime: expireTime,
	}
}

func (r EmailChangeRequest) Email() Email {
	return r.email
}

func (r EmailChangeRequest) Token() EmailChangeRequestToken {
	return r.token
}

func (r EmailChangeRequest) CreateTime() time.Time {
	return r.createTime
}

func (r EmailChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}

func (r EmailChangeRequest) Equal(other EmailChangeRequest) bool {
	return r.email.Equal(other.email) && r.createTime.Equal(other.createTime) && r.expireTime.Equal(other.expireTime)
}
