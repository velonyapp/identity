package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[*event.SessionCreated] = (*SessionCreatedHandler)(nil)

type SessionCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewSessionCreatedHandler(
	eventPublisher port.EventPublisher,
) *SessionCreatedHandler {
	return &SessionCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *SessionCreatedHandler) Execute(ctx context.Context, domainEvent *event.SessionCreated) error {
	integrationEvent := integrationevent.NewSessionCreated(
		domainEvent.AggregateID(),
		domainEvent.UserID().Value(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
