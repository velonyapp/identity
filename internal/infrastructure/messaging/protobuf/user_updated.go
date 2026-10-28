package protobuf

import (
	v1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userUpdatedPayload(event integrationevent.UserUpdated) *v1.UserUpdatedPayload {
	return &v1.UserUpdatedPayload{
		UserId:   event.AggregateID(),
		Username: event.Username(),
		Email:    event.Email(),
		Version:  event.Version(),
	}
}
