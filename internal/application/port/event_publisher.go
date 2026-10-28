package port

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
)

type EventPublisher interface {
	Publish(ctx context.Context, event integrationevent.IntegrationEvent) error
}
