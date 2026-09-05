package port

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
)

type OutboxMessage struct {
	PartitionKey string
	Event        *v1.Event
}

type OutboxPublisher interface {
	PublishMessage(ctx context.Context, message OutboxMessage) error
	PublishMessages(ctx context.Context, messages []OutboxMessage) error
}
