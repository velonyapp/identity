package protobuf

import (
	v1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userCreatedPayload(event integrationevent.UserCreated) *v1.UserCreatedPayload {
	return &v1.UserCreatedPayload{
		UserId:   event.AggregateID(),
		Username: event.Username(),
		Version:  event.Version(),
	}
}
