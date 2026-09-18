package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/infrastructure/messaging/event"

	"google.golang.org/protobuf/encoding/protojson"
)

type outboxPublisher struct {
	db      *sql.DB
	encoder *event.Encoder
}

func NewOutboxPublisher(
	db *sql.DB,
	encoder *event.Encoder,
) port.OutboxPublisher {
	return &outboxPublisher{
		db:      db,
		encoder: encoder,
	}
}

func (pub *outboxPublisher) PublishMessage(ctx context.Context, message port.OutboxMessage) error {
	event, err := pub.encoder.Encode(message.Event)
	if err != nil {
		return err
	}

	payload, err := protojson.Marshal(event)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO outbox_messages (
			id,
			type,
			aggregate_id,
			aggregate_type,
			occur_time,
			payload,
			partition_key
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		event.Id,
		event.Type,
		event.AggregateId,
		event.AggregateType,
		event.OccurTime,
		string(payload),
		message.PartitionKey,
	)

	return err
}

func (pub *outboxPublisher) PublishMessages(ctx context.Context, messages []port.OutboxMessage) error {
	if len(messages) == 0 {
		return nil
	}

	const prefix = `
		INSERT INTO outbox_messages (
			id,
			type,
			aggregate_id,
			aggregate_type,
			occur_time,
			payload,
			partition_key
		)
		VALUES
	`

	var query strings.Builder
	query.WriteString(prefix)

	args := make([]any, 0, len(messages)*7)

	for i, message := range messages {
		event, err := pub.encoder.Encode(message.Event)
		if err != nil {
			return err
		}

		payload, err := protojson.Marshal(event)
		if err != nil {
			return err
		}

		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString("(?, ?, ?, ?, ?, ?, ?)")

		args = append(args,
			event.Id,
			event.Type,
			event.AggregateId,
			event.AggregateType,
			event.OccurTime,
			string(payload),
			message.PartitionKey,
		)
	}

	_, err := executor(ctx, pub.db).ExecContext(ctx, query.String(), args...)

	return err
}
