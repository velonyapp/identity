package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func sessionRefreshedPayload(event integrationevent.SessionRefreshed) *eventv1.SessionRefreshedPayload {
	return &eventv1.SessionRefreshedPayload{
		RefreshTime: timestamppb.New(event.RefreshTime),
	}
}
