package vo

import "time"

type AvatarChangeRequest struct {
	avatarID   AvatarID
	expireTime time.Time
}

func NewAvatarChangeRequest(
	avatarID AvatarID,
	expireTime time.Time,
) AvatarChangeRequest {
	return AvatarChangeRequest{
		avatarID:   avatarID,
		expireTime: expireTime,
	}
}

func (r AvatarChangeRequest) AvatarID() AvatarID {
	return r.avatarID
}

func (r AvatarChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}
