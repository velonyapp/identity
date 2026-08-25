package domainevent

import (
	"context"

	v1 "github.com/velony-app/identity/gen/event/v1"
	integrationevent "github.com/velony-app/identity/internal/application/event"
	"github.com/velony-app/identity/internal/application/port"
	domainevent "github.com/velony-app/identity/internal/domain/event"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func (h *UserUsernameChangedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserUsernameChanged,
) error {
	payload := &v1.UserUsernameChangedPayload{
		UserId:     domainEvent.AggregateID(),
		Username:   domainEvent.Username.Value(),
		UpdateTime: timestamppb.New(domainEvent.UpdateTime.Value()),
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		"user.username.changed.v1",
		payload,
	)
	if err != nil {
		return err
	}

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        integrationEvent,
		},
	)
}
