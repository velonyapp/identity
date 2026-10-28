package mysql

import (
	"context"
	"database/sql"

	"github.com/velonyapp/identity/internal/application/job"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/infrastructure/messaging/protobuf"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ port.JobPublisher = (*jobPublisher)(nil)

type jobPublisher struct {
	db      *sql.DB
	encoder *protobuf.JobEncoder
}

func NewJobPublisher(
	db *sql.DB,
	encoder *protobuf.JobEncoder,
) port.JobPublisher {
	return &jobPublisher{
		db:      db,
		encoder: encoder,
	}
}

func (pub *jobPublisher) Publish(ctx context.Context, job job.Job) error {
	jobPayload, err := pub.encoder.Encode(job)
	if err != nil {
		return err
	}

	payload, err := protojson.Marshal(jobPayload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO outbox_jobs (
			id,
			type,
			payload,
			enqueue_time
		)
		VALUES (?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		job.ID(),
		job.Type(),
		string(payload),
		job.EnqueueTime(),
	)

	return err
}
