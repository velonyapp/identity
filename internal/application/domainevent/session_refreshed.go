package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type SessionRefreshedHandler struct {
	eventPublisher port.EventPublisher
}

func NewSessionRefreshedHandler(
	eventPublisher port.EventPublisher,
) *SessionRefreshedHandler {
	return &SessionRefreshedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *SessionRefreshedHandler) Execute(ctx context.Context, domainEvent event.SessionRefreshed) error {
	integrationEvent := integrationevent.NewSessionRefreshed(
		domainEvent.AggregateID(),
		domainEvent.RefreshTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
