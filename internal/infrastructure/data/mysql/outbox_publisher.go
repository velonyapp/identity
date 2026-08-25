package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/velony-app/identity/internal/application/port"

	"google.golang.org/protobuf/encoding/protojson"
)

type outboxPublisher struct {
	db *sql.DB
}

func NewOutboxPublisher(db *sql.DB) port.OutboxPublisher {
	return &outboxPublisher{db: db}
}

func (pub *outboxPublisher) PublishMessage(
	ctx context.Context,
	message port.OutboxMessage,
) error {
	payload, err := protojson.Marshal(message.Event.Payload)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO outbox_messages (
			id,
			partition_key,
			source,
			type,
			occur_time,
			payload
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		message.Event.Id,
		message.PartitionKey,
		message.Event.Source,
		message.Event.Type,
		message.Event.OccurTime.AsTime(),
		string(payload),
	)

	return err
}

func (pub *outboxPublisher) PublishMessages(
	ctx context.Context,
	messages []port.OutboxMessage,
) error {
	if len(messages) == 0 {
		return nil
	}

	const prefix = `
		INSERT INTO outbox_messages (
			id,
			partition_key,
			source,
			type,
			occur_time,
			payload
		)
		VALUES
	`

	var query strings.Builder
	query.WriteString(prefix)

	args := make([]any, 0, len(messages)*6)

	for i, message := range messages {
		if message.Event == nil {
			return errors.New("outbox message event is nil")
		}

		payload, err := protojson.Marshal(message.Event.Payload)
		if err != nil {
			return err
		}

		if i > 0 {
			query.WriteString(",")
		}

		query.WriteString("(?, ?, ?, ?, ?, ?)")

		args = append(args,
			message.Event.Id,
			message.PartitionKey,
			message.Event.Source,
			message.Event.Type,
			message.Event.OccurTime.AsTime(),
			string(payload),
		)
	}

	_, err := executor(ctx, pub.db).ExecContext(ctx, query.String(), args...)

	return err
}
