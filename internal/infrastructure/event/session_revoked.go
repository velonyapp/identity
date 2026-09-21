package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func sessionRevokedPayload(event integrationevent.SessionRevoked) *eventv1.SessionRevokedPayload {
	return &eventv1.SessionRevokedPayload{
		RevokeTime: timestamppb.New(event.RevokeTime),
	}
}
