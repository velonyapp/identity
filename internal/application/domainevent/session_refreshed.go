package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type SessionRefreshedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewSessionRefreshedHandler(
	outboxPublisher port.OutboxPublisher,
) *SessionRefreshedHandler {
	return &SessionRefreshedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *SessionRefreshedHandler) Execute(ctx context.Context, domainEvent event.SessionRefreshed) error {
	integrationEvent := integrationevent.NewSessionRefreshed(
		domainEvent.AggregateID(),
		domainEvent.RefreshTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
