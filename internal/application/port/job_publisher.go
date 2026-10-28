package port

import (
	"context"

	"github.com/velonyapp/identity/internal/application/job"
)

type JobPublisher interface {
	Publish(ctx context.Context, job job.Job) error
}
