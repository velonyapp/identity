package port

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
)

type OutboxMessage struct {
	PartitionKey string
	Event        integrationevent.IntegrationEvent
}

type OutboxPublisher interface {
	PublishMessage(ctx context.Context, message OutboxMessage) error
	PublishMessages(ctx context.Context, messages []OutboxMessage) error
}
