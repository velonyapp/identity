package protobuf

import (
	v1 "github.com/velonyapp/asset/gen/job/v1"
	"github.com/velonyapp/identity/internal/application/job"
)

func removeAvatarPayload(job job.RemoveAvatar) *v1.DeleteImagePayload {
	return &v1.DeleteImagePayload{
		ImageId: job.AvatarID,
	}
}
