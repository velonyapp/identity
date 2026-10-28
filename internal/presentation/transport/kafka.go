package transport

import (
	"context"
	"errors"
	"slices"
	"strings"

	assetv1 "github.com/velonyapp/asset/gen/event/v1"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/presentation/service"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v3/transport"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

const (
	groupID = "identity"

	assetImageTopic = "asset.image"

	assetImageActiveEvent = "asset.image.active"
)

var _ transport.Server = (*KafkaConsumer)(nil)

type KafkaConsumer struct {
	c   *conf.Transport
	svc *service.Consumer

	consumer *kafka.Consumer
}

func NewKafkaConsumer(
	c *conf.Transport,
	svc *service.Consumer,
) *KafkaConsumer {
	return &KafkaConsumer{
		c:   c,
		svc: svc,
	}
}

func (k *KafkaConsumer) handleMessage(ctx context.Context, message *kafka.Message) error {
	ctx, cancel := context.WithTimeout(ctx, k.c.Kafka.Timeout.AsDuration())
	defer cancel()

	var eventType string
	for _, header := range message.Headers {
		switch header.Key {
		case "type":
			eventType = string(header.Value)
		}
	}

	switch eventType {
	case assetImageActiveEvent:
		event := new(assetv1.ImageActivePayload)
		if err := protojson.Unmarshal(message.Value, event); err != nil {
			return ErrInvalidMessage
		}

		if !slices.Contains(event.Tags, "identity") || !slices.Contains(event.Tags, "user") {
			return nil
		}

		userID, ok := event.Metadata["user_id"]
		if !ok {
			return ErrInvalidMessage
		}

		if err := k.svc.ConfirmUserAvatarChange(ctx, service.ConfirmUserAvatarChangeRequest{
			UserID: userID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaConsumer) Start(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(k.c.Kafka.Brokers, ","),
		"group.id":                 groupID,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
	})
	if err != nil {
		return err
	}

	topics := []string{
		assetImageTopic,
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		return err
	}

	for !consumer.IsClosed() {
		event := consumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			if err := k.handleMessage(ctx, e); err != nil && !errors.Is(err, ErrInvalidMessage) {
				break
			}

			consumer.StoreMessage(e)

		case kafka.Error:
			if e.IsFatal() {
				return errors.New(e.Error())
			}
		}
	}

	return nil
}

func (k *KafkaConsumer) Stop(ctx context.Context) error {
	return k.consumer.Close()
}
