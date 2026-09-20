package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type SessionCreatedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewSessionCreatedHandler(
	outboxPublisher port.OutboxPublisher,
) *SessionCreatedHandler {
	return &SessionCreatedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *SessionCreatedHandler) Execute(ctx context.Context, domainEvent event.SessionCreated) error {
	integrationEvent := integrationevent.NewSessionCreated(
		domainEvent.AggregateID(),
		domainEvent.UserID.Value(),
		domainEvent.CreateTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
