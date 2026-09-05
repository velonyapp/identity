package event

import (
	v1 "github.com/velonyapp/identity/gen/event/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewIntegrationEvent(
	eventType string,
	payload proto.Message,
) (*v1.Event, error) {
	packed, err := anypb.New(payload)
	if err != nil {
		return nil, err
	}

	return &v1.Event{
		Id:        uuid.NewString(),
		Source:    "velony.identity",
		Type:      eventType,
		OccurTime: timestamppb.Now(),
		Payload:   packed,
	}, nil
}
