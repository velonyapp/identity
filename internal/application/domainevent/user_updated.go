package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[event.UserUpdated] = (*UserUpdatedHandler)(nil)

type UserUpdatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserUpdatedHandler(
	eventPublisher port.EventPublisher,
) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserUpdatedHandler) Execute(ctx context.Context, domainEvent event.UserUpdated) error {
	var emailValue *string
	if email := domainEvent.Email(); email != nil {
		value := email.String()
		emailValue = &value
	}

	integrationEvent := integrationevent.NewUserUpdated(
		domainEvent.AggregateID(),
		domainEvent.UpdateTime(),
		domainEvent.Username().String(),
		emailValue,
		domainEvent.Version(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
