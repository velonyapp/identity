package job

import "time"

var _ Job = (*RemoveAvatar)(nil)

type RemoveAvatar struct {
	BaseJob

	AvatarID string
}

func NewRemoveAvatar(
	AvatarID string,
	enqueueTime time.Time,
) RemoveAvatar {
	return RemoveAvatar{
		BaseJob: NewBaseJob(enqueueTime),

		AvatarID: AvatarID,
	}
}

func (e RemoveAvatar) Type() string {
	return "asset.image.delete"
}

func (e RemoveAvatar) ImageID() string {
	return e.AvatarID
}
