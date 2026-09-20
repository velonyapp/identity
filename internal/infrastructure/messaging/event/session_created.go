package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func sessionCreatedPayload(event integrationevent.SessionCreated) *eventv1.SessionCreatedPayload {
	return &eventv1.SessionCreatedPayload{
		UserId:     event.UserID,
		CreateTime: timestamppb.New(event.CreateTime),
	}
}
