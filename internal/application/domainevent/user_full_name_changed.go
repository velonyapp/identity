package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserFullNameChangedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserFullNameChangedHandler(
	eventPublisher port.EventPublisher,
) *UserFullNameChangedHandler {
	return &UserFullNameChangedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserFullNameChangedHandler) Execute(ctx context.Context, domainEvent event.UserFullNameChanged) error {
	integrationEvent := integrationevent.NewUserFullNameChanged(
		domainEvent.AggregateID(),
		domainEvent.OldFullName.Value(),
		domainEvent.NewFullName.Value(),
		domainEvent.UpdateTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
