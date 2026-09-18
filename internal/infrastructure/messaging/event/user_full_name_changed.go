package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userFullNameChangedPayload(event integrationevent.UserFullNameChanged) *eventv1.UserFullNameChangedPayload {
	return &eventv1.UserFullNameChangedPayload{
		OldFullName: event.OldFullName,
		NewFullName: event.NewFullName,
		UpdateTime:  timestamppb.New(event.UpdateTime),
	}
}
