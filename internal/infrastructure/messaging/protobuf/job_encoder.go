package protobuf

import (
	"fmt"

	"github.com/velonyapp/identity/internal/application/job"

	"google.golang.org/protobuf/proto"
)

type JobEncoder struct{}

func NewJobEncoder() *JobEncoder {
	return &JobEncoder{}
}

func (e *JobEncoder) Encode(jobPayload job.Job) (proto.Message, error) {
	switch jobPayload := jobPayload.(type) {
	case job.SendEmail:
		return sendEmailPayload(jobPayload), nil
	case job.RemoveAvatar:
		return removeAvatarPayload(jobPayload), nil
	default:
		return nil, fmt.Errorf("unknown job %T", jobPayload)
	}
}
