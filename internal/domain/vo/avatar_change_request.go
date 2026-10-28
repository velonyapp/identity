package vo

import "time"

type AvatarChangeRequest struct {
	avatarID   AvatarID
	createTime time.Time
	expireTime time.Time
}

func NewAvatarChangeRequest(
	avatarID AvatarID,
	createTime time.Time,
	expireTime time.Time,
) AvatarChangeRequest {
	return AvatarChangeRequest{
		avatarID:   avatarID,
		createTime: createTime,
		expireTime: expireTime,
	}
}

func (r AvatarChangeRequest) AvatarID() AvatarID {
	return r.avatarID
}

func (r AvatarChangeRequest) CreateTime() time.Time {
	return r.createTime
}

func (r AvatarChangeRequest) ExpireTime() time.Time {
	return r.expireTime
}

func (r AvatarChangeRequest) Equal(other AvatarChangeRequest) bool {
	return r.avatarID.Equal(other.avatarID) && r.createTime.Equal(other.createTime) && r.expireTime.Equal(other.expireTime)
}
