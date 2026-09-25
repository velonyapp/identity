package event

import (
	"fmt"

	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(event integrationevent.IntegrationEvent) (*eventv1.Event, error) {
	if event == nil {
		return nil, fmt.Errorf("integration event is nil")
	}

	payload, err := e.payload(event)
	if err != nil {
		return nil, err
	}

	payloadAny, err := anypb.New(payload)
	if err != nil {
		return nil, err
	}

	return &eventv1.Event{
		Id:            event.ID(),
		Type:          event.Type(),
		AggregateId:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		OccurTime:     timestamppb.New(event.OccurTime()),
		Payload:       payloadAny,
	}, nil
}

func (e *Encoder) payload(event integrationevent.IntegrationEvent) (proto.Message, error) {
	switch event := event.(type) {

	case *integrationevent.UserCreated:
		if event == nil {
			return nil, fmt.Errorf("UserCreated event is nil")
		}

		return userCreatedPayload(*event), nil

	case *integrationevent.UserUsernameChanged:
		if event == nil {
			return nil, fmt.Errorf("UserUsernameChanged event is nil")
		}

		return userUsernameChangedPayload(*event), nil

	case *integrationevent.UserEmailChanged:
		if event == nil {
			return nil, fmt.Errorf("UserEmailChanged event is nil")
		}

		return userEmailChangedPayload(*event), nil

	case *integrationevent.UserFullNameChanged:
		if event == nil {
			return nil, fmt.Errorf("UserFullNameChanged event is nil")
		}

		return userFullNameChangedPayload(*event), nil

	case *integrationevent.UserAvatarChanged:
		if event == nil {
			return nil, fmt.Errorf("UserAvatarKeyChanged event is nil")
		}

		return userAvatarKeyChangedPayload(*event), nil

	case *integrationevent.UserDeleted:
		if event == nil {
			return nil, fmt.Errorf("UserDeleted event is nil")
		}

		return userDeletedPayload(*event), nil

	case *integrationevent.SessionCreated:
		if event == nil {
			return nil, fmt.Errorf("SessionCreated event is nil")
		}

		return sessionCreatedPayload(*event), nil

	case *integrationevent.SessionRefreshed:
		if event == nil {
			return nil, fmt.Errorf("SessionRefreshed event is nil")
		}

		return sessionRefreshedPayload(*event), nil

	case *integrationevent.SessionRevoked:
		if event == nil {
			return nil, fmt.Errorf("SessionRevoked event is nil")
		}

		return sessionRevokedPayload(*event), nil

	default:
		return nil, fmt.Errorf("unknown integration event %T", event)
	}
}
