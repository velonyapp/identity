package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userDeletedPayload(event integrationevent.UserDeleted) *eventv1.UserDeletedPayload {
	return &eventv1.UserDeletedPayload{
		DeleteTime: timestamppb.New(event.DeleteTime),
	}
}
