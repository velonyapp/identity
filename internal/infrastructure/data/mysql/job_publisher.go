
package mysql

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/velonyapp/identity/internal/application/job"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/infrastructure/messaging/protobuf"

	"google.golang.org/protobuf/encoding/protojson"
)

var _ port.JobPublisher = (*jobPublisher)(nil)

type jobPublisher struct {
	db      *sql.DB
	encoder *protobuf.JobEncoder
	log     *slog.Logger
}

func NewJobPublisher(
	db *sql.DB,
	encoder *protobuf.JobEncoder,
	log *slog.Logger,
) port.JobPublisher {
	return &jobPublisher{
		db:      db,
		encoder: encoder,
		log: log.With(
			slog.String("component", "mysql.job_publisher"),
		),
	}
}

func (pub *jobPublisher) Publish(ctx context.Context, job job.Job) error {
	jobPayload, err := pub.encoder.Encode(job)
	if err != nil {
		pub.log.Error(
			"failed to encode job",
			slog.Any("job_id", job.ID()),
			slog.Any("job_type", job.Type()),
			slog.Any("error", err),
		)
		return err
	}

	payload, err := protojson.Marshal(jobPayload)
	if err != nil {
		pub.log.Error(
			"failed to marshal job payload",
			slog.Any("job_id", job.ID()),
			slog.Any("job_type", job.Type()),
			slog.Any("error", err),
		)
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
	if err != nil {
		pub.log.Error(
			"failed to insert job into outbox",
			slog.Any("job_id", job.ID()),
			slog.Any("job_type", job.Type()),
			slog.Any("error", err),
		)
	}

	return err
}
