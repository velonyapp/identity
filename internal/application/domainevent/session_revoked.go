package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[*event.SessionRevoked] = (*SessionRevokedHandler)(nil)

type SessionRevokedHandler struct {
	eventPublisher port.EventPublisher
}

func NewSessionRevokedHandler(
	eventPublisher port.EventPublisher,
) *SessionRevokedHandler {
	return &SessionRevokedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *SessionRevokedHandler) Execute(ctx context.Context, domainEvent *event.SessionRevoked) error {
	integrationEvent := integrationevent.NewSessionRevoked(
		domainEvent.AggregateID(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
