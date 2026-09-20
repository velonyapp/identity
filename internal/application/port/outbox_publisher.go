package port

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
)

type OutboxPublisher interface {
	PublishMessage(ctx context.Context, message integrationevent.IntegrationEvent) error
	PublishMessages(ctx context.Context, messages []integrationevent.IntegrationEvent) error
}
