package protobuf

import (
	v1 "github.com/velonyapp/email/gen/job/v1"
	"github.com/velonyapp/identity/internal/application/job"
)

func sendEmailPayload(job job.SendEmail) *v1.SendEmailPayload {
	return &v1.SendEmailPayload{
		From:           "no-reply@velony.app",
		To:             []string{job.To()},
		Template:       job.Template(),
		Variables:      job.Variables(),
		IdempotencyKey: job.ID(),
	}
}
