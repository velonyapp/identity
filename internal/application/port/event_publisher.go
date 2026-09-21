package port

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
)

type EventPublisher interface {
	Publish(ctx context.Context, integrationEvent integrationevent.IntegrationEvent) error
	PublishBatch(ctx context.Context, integrationEvents []integrationevent.IntegrationEvent) error
}
