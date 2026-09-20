package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type SessionRevokedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewSessionRevokedHandler(
	outboxPublisher port.OutboxPublisher,
) *SessionRevokedHandler {
	return &SessionRevokedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *SessionRevokedHandler) Execute(ctx context.Context, domainEvent event.SessionRevoked) error {
	integrationEvent := integrationevent.NewSessionRevoked(
		domainEvent.AggregateID(),
		domainEvent.RevokeTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
