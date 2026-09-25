package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserCreatedHandler(
	eventPublisher port.EventPublisher,
) *UserCreatedHandler {
	return &UserCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserCreatedHandler) Execute(ctx context.Context, domainEvent *event.UserCreated) error {
	integrationEvent := integrationevent.NewUserCreated(
		domainEvent.AggregateID(),
		domainEvent.Username().Value(),
		domainEvent.FullName().Value(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
