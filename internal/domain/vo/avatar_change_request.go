package vo

import "time"

type AvatarChangeRequest struct {
	value      AvatarKey
	time       time.Time
	expireTime time.Time
}

func NewAvatarChangeRequest(
	value AvatarKey,
	time time.Time,
	expireTime time.Time,
) AvatarChangeRequest {
	return AvatarChangeRequest{
		value:      value,
		time:       time,
		expireTime: expireTime,
	}
}

func (r AvatarChangeRequest) Value() AvatarKey {
	return r.value
}

func (r AvatarChangeRequest) Time() time.Time {
	return r.time
}

func (r AvatarChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}
