package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserEmailChangedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserEmailChangedHandler(
	eventPublisher port.EventPublisher,
) *UserEmailChangedHandler {
	return &UserEmailChangedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserEmailChangedHandler) Execute(
	ctx context.Context,
	domainEvent event.UserEmailChanged,
) error {
	var oldEmail, newEmail *string

	if domainEvent.OldEmail != nil {
		value := domainEvent.OldEmail.String()
		oldEmail = &value
	}
	if domainEvent.NewEmail != nil {
		value := domainEvent.NewEmail.String()
		newEmail = &value
	}

	integrationEvent := integrationevent.NewUserEmailChanged(
		domainEvent.AggregateID(),
		oldEmail,
		newEmail,
		domainEvent.UpdateTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
