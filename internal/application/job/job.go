package job

import (
	"time"

	"github.com/google/uuid"
)

type Job interface {
	ID() string
	Type() string
	EnqueueTime() time.Time
}

type BaseJob struct {
	id          string
	enqueueTime time.Time
}

func NewBaseJob(enqueueTime time.Time) BaseJob {
	return BaseJob{
		id:          uuid.Must(uuid.NewV7()).String(),
		enqueueTime: enqueueTime,
	}
}

func (e BaseJob) ID() string             { return e.id }
func (e BaseJob) EnqueueTime() time.Time { return e.enqueueTime }
