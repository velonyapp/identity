package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserDeletedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserDeletedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserDeletedHandler {
	return &UserDeletedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserDeletedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserDeleted,
) error {
	payload := &v1.UserDeletedPayload{
		UserId:     domainEvent.AggregateID(),
		DeleteTime: timestamppb.New(domainEvent.DeleteTime.Value()),
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		"user.deleted.v1",
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
