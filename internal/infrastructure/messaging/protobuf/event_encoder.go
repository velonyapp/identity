package protobuf

import (
	"fmt"

	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
)

type EventEncoder struct{}

func NewEventEncoder() *EventEncoder {
	return &EventEncoder{}
}

func (e *EventEncoder) Encode(eventPayload integrationevent.IntegrationEvent) (proto.Message, error) {
	switch eventPayload := eventPayload.(type) {
	case integrationevent.UserCreated:
		return userCreatedPayload(eventPayload), nil
	case integrationevent.UserUpdated:
		return userUpdatedPayload(eventPayload), nil
	case integrationevent.UserDeleted:
		return userDeletedPayload(eventPayload), nil
	default:
		return nil, fmt.Errorf("unknown integration event %T", eventPayload)
	}
}
