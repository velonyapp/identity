package transport

import (
	"context"
	"errors"
	"strings"
	"sync"

	assetv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/command"
	"github.com/velonyapp/identity/internal/application/common"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/domain/entity"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v3/transport"
	"google.golang.org/protobuf/proto"
)

var _ transport.Server = (*KafkaConsumer)(nil)

type KafkaConsumer struct {
	c *conf.Transport

	reconstituteAvatarHandler      *command.ReconstituteAvatarHandler
	deleteAvatarHandler            *command.DeleteAvatarHandler
	confirmUserAvatarChangeHandler *command.ConfirmUserAvatarChangeHandler

	consumers map[string]*kafka.Consumer

	wg sync.WaitGroup
}

func NewKafkaConsumer(
	c *conf.Transport,
	reconstituteAvatarHandler *command.ReconstituteAvatarHandler,
	deleteAvatarHandler *command.DeleteAvatarHandler,
	confirmUserAvatarChangeHandler *command.ConfirmUserAvatarChangeHandler,
) *KafkaConsumer {
	return &KafkaConsumer{
		c:                              c,
		reconstituteAvatarHandler:      reconstituteAvatarHandler,
		deleteAvatarHandler:            deleteAvatarHandler,
		confirmUserAvatarChangeHandler: confirmUserAvatarChangeHandler,

		consumers: make(map[string]*kafka.Consumer),
	}
}

func (kc *KafkaConsumer) registerAllConsumers() error {
	var errs []error

	if err := kc.registerConsumer(
		kc.c.Kafka.Consumers.AssetImage.Topic,
		kc.c.Kafka.Consumers.AssetImage.GroupId,
	); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (kc *KafkaConsumer) handleMessage(ctx context.Context, message *kafka.Message) error {
	var eventType string
	for _, header := range message.Headers {
		if header.Key == "type" {
			eventType = string(header.Value)
			break
		}
	}

	switch eventType {
	case kc.c.Kafka.Consumers.AssetImage.Events.Created:
		var event assetv1.Event
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			return err
		}

		var payload assetv1.ImageCreatedPayload
		if err := event.GetPayload().UnmarshalTo(&payload); err != nil {
			return err
		}

		if _, err := kc.reconstituteAvatarHandler.Execute(ctx, &command.ReconstituteAvatar{
			AvatarID: event.AggregateId,
			Key:      payload.ObjectKey,
		}); err != nil {
			return err
		}

	case kc.c.Kafka.Consumers.AssetImage.Events.ObjectExistenceUpdated:
		var event assetv1.Event
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			return err
		}

		var payload assetv1.ImageObjectExistenceUpdatedPayload
		if err := event.GetPayload().UnmarshalTo(&payload); err != nil {
			return err
		}

		if !payload.ObjectExists {
			break
		}

		if _, err := kc.confirmUserAvatarChangeHandler.Execute(ctx, &command.ConfirmUserAvatarChange{
			AvatarID: event.AggregateId,
		}); err != nil {
			switch {
			case errors.Is(err, common.ErrUserNotFound),
				errors.Is(err, entity.ErrUserDeleted):
			default:
				return err
			}
		}

	case kc.c.Kafka.Consumers.AssetImage.Events.Deleted:
		var event assetv1.Event
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			return err
		}

		var payload assetv1.ImageDeletedPayload
		if err := event.GetPayload().UnmarshalTo(&payload); err != nil {
			return err
		}

		kc.deleteAvatarHandler.Execute(ctx, &command.DeleteAvatar{
			AvatarID: event.AggregateId,
		})
	}

	return nil
}

func (kc *KafkaConsumer) registerConsumer(topic, groupID string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(kc.c.Kafka.Brokers, ","),
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	})
	if err != nil {
		return err
	}

	if err := consumer.Subscribe(topic, nil); err != nil {
		consumer.Close()
		return err
	}

	kc.consumers[topic] = consumer

	return nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	if err := kc.registerAllConsumers(); err != nil {
		return err
	}

	errChan := make(chan error, len(kc.consumers))

	for _, consumer := range kc.consumers {
		kc.wg.Go(func() {
			for !consumer.IsClosed() {
				event := consumer.Poll(int(kc.c.Kafka.Timeout.GetSeconds()) * 1000)
				if event == nil {
					continue
				}

				switch e := event.(type) {
				case *kafka.Message:
					if err := kc.handleMessage(ctx, e); err != nil {
						break
					}

					consumer.StoreMessage(e)

				case kafka.Error:
					if e.IsFatal() {
						errChan <- e
						return
					}
				}
			}
		})
	}

	select {
	case <-ctx.Done():
		return nil

	case err := <-errChan:
		return err
	}
}

func (kc *KafkaConsumer) Stop(ctx context.Context) error {
	var errs []error

	for topic, consumer := range kc.consumers {
		if consumer == nil {
			delete(kc.consumers, topic)
			continue
		}

		if err := consumer.Close(); err != nil {
			errs = append(errs, err)
		}

		delete(kc.consumers, topic)
	}

	kc.wg.Wait()

	return nil
}
