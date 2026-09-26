package vo

import "time"

type AvatarChangeRequest struct {
	avatarKey  AvatarKey
	time       time.Time
	expireTime time.Time
}

func NewAvatarChangeRequest(
	avatarKey AvatarKey,
	time time.Time,
	expireTime time.Time,
) AvatarChangeRequest {
	return AvatarChangeRequest{
		avatarKey:  avatarKey,
		time:       time,
		expireTime: expireTime,
	}
}

func (r AvatarChangeRequest) AvatarKey() AvatarKey {
	return r.avatarKey
}

func (r AvatarChangeRequest) Time() time.Time {
	return r.time
}

func (r AvatarChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}
