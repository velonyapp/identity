package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[*event.UserDeleted] = (*UserDeletedHandler)(nil)

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

func (h *UserDeletedHandler) Execute(ctx context.Context, domainEvent *event.UserDeleted) error {
	integrationEvent := integrationevent.NewUserDeleted(
		domainEvent.AggregateID(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
