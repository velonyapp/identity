package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserEmailChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserEmailChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserEmailChangedHandler {
	return &UserEmailChangedHandler{
		outboxPublisher: outboxPublisher,
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

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
