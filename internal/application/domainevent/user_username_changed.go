package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserUsernameChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserUsernameChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserUsernameChangedHandler {
	return &UserUsernameChangedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserUsernameChangedHandler) Execute(ctx context.Context, domainEvent event.UserUsernameChanged) error {
	integrationEvent := integrationevent.NewUserUsernameChanged(
		domainEvent.AggregateID(),
		domainEvent.OldUsername.Value(),
		domainEvent.NewUsername.Value(),
		domainEvent.UpdateTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
