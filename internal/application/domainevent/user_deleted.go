package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserDeletedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserDeletedHandler(
	eventPublisher port.EventPublisher,
) *UserDeletedHandler {
	return &UserDeletedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserDeletedHandler) Execute(ctx context.Context, domainEvent event.UserDeleted) error {
	integrationEvent := integrationevent.NewUserDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
