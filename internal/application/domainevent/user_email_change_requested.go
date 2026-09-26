package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[*event.UserEmailChangeRequested] = (*UserEmailChangeRequestedHandler)(nil)

type UserEmailChangeRequestedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserEmailChangeRequestedHandler(
	eventPublisher port.EventPublisher,
) *UserEmailChangeRequestedHandler {
	return &UserEmailChangeRequestedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserEmailChangeRequestedHandler) Execute(ctx context.Context, domainEvent *event.UserEmailChangeRequested) error {
	integrationEvent := integrationevent.NewUserEmailChangeRequested(
		domainEvent.AggregateID(),
		domainEvent.Email().Value(),
		domainEvent.ExpireTime(),
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
