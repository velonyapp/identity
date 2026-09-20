package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/velonyapp/identity/internal/application/integrationevent"
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

func (pub *outboxPublisher) PublishMessage(ctx context.Context, integrationEvent integrationevent.IntegrationEvent) error {
	event, err := pub.encoder.Encode(integrationEvent)
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
			payload
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = executor(ctx, pub.db).ExecContext(ctx, query,
		event.Id,
		event.Type,
		event.AggregateId,
		event.AggregateType,
		event.OccurTime.AsTime(),
		string(payload),
	)

	return err
}

func (pub *outboxPublisher) PublishMessages(ctx context.Context, integrationEvents []integrationevent.IntegrationEvent) error {
	if len(integrationEvents) == 0 {
		return nil
	}

	const prefix = `
		INSERT INTO outbox_messages (
			id,
			type,
			aggregate_id,
			aggregate_type,
			occur_time,
			payload
		)
		VALUES
	`

	var query strings.Builder
	query.WriteString(prefix)

	args := make([]any, 0, len(integrationEvents)*6)

	for i, integrationEvent := range integrationEvents {
		event, err := pub.encoder.Encode(integrationEvent)
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
		query.WriteString("(?, ?, ?, ?, ?, ?)")

		args = append(args,
			event.Id,
			event.Type,
			event.AggregateId,
			event.AggregateType,
			event.OccurTime.AsTime(),
			string(payload),
		)
	}

	_, err := executor(ctx, pub.db).ExecContext(ctx, query.String(), args...)

	return err
}
