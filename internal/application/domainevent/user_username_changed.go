package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserUsernameChangedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserUsernameChangedHandler(
	eventPublisher port.EventPublisher,
) *UserUsernameChangedHandler {
	return &UserUsernameChangedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserUsernameChangedHandler) Execute(ctx context.Context, domainEvent *event.UserUsernameChanged) error {
	integrationEvent := integrationevent.NewUserUsernameChanged(
		domainEvent.AggregateID(),
		domainEvent.OldUsername().Value(),
		domainEvent.NewUsername().Value(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
