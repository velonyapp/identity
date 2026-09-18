package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userUsernameChangedPayload(event integrationevent.UserUsernameChanged) *eventv1.UserUsernameChangedPayload {
	return &eventv1.UserUsernameChangedPayload{
		OldUsername: event.OldUsername,
		NewUsername: event.NewUsername,
		UpdateTime:  timestamppb.New(event.UpdateTime),
	}
}
