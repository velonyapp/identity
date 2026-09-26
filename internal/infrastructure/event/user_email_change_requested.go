package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userEmailChangeRequestedPayload(event integrationevent.UserEmailChangeRequested) *eventv1.UserEmailChangeRequestedPayload {
	return &eventv1.UserEmailChangeRequestedPayload{
		Email:      event.Email(),
		ExpireTime: timestamppb.New(event.ExpireTime()),
	}
}
