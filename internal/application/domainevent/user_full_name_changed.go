package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func (h *UserFullNameChangedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserFullNameChanged,
) error {
	payload := &v1.UserFullNameChangedPayload{
		UserId:     domainEvent.AggregateID(),
		FullName:   domainEvent.FullName.Value(),
		UpdateTime: timestamppb.New(domainEvent.UpdateTime.Value()),
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		"user.full-name.changed.v1",
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
