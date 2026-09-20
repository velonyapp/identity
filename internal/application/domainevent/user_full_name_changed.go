package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserFullNameChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserFullNameChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserFullNameChangedHandler {
	return &UserFullNameChangedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserFullNameChangedHandler) Execute(ctx context.Context, domainEvent event.UserFullNameChanged) error {
	integrationEvent := integrationevent.NewUserFullNameChanged(
		domainEvent.AggregateID(),
		domainEvent.OldFullName.Value(),
		domainEvent.NewFullName.Value(),
		domainEvent.UpdateTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
