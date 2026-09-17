package event

import (
	"time"

	v1 "github.com/velonyapp/identity/gen/event/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO:
// Move protobuf-specific integration event construction out of
// the application layer into an anti-corruption layer.
func NewIntegrationEvent(
	eventType string,
	aggregateID string,
	aggregateType string,
	occurTime time.Time,
	payload proto.Message,
) (*v1.Event, error) {
	packed, err := anypb.New(payload)
	if err != nil {
		return nil, err
	}

	return &v1.Event{
		Id:            uuid.NewString(),
		Type:          eventType,
		AggregateId:   aggregateID,
		AggregateType: aggregateType,
		OccurTime:     timestamppb.New(occurTime),
		Payload:       packed,
	}, nil
}
