package transport

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

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
	log *slog.Logger

	consumer *kafka.Consumer
}

func NewKafkaConsumer(
	c *conf.Transport,
	svc *service.Consumer,
	log *slog.Logger,
) *KafkaConsumer {
	return &KafkaConsumer{
		c:   c,
		svc: svc,
		log: log,
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
		k.log.ErrorContext(ctx,
			"failed to create Kafka consumer",
			slog.Any("error", err),
		)
		return err
	}

	k.consumer = consumer

	topics := []string{
		assetImageTopic,
	}

	if err := consumer.SubscribeTopics(topics, nil); err != nil {
		k.log.ErrorContext(ctx,
			"failed to subscribe to Kafka topics",
			slog.Any("topics", topics),
			slog.Any("error", err),
		)
		return err
	}

	k.log.InfoContext(ctx,
		"Kafka consumer subscribed to topics",
		slog.String("group_id", groupID),
		slog.Any("topics", topics),
	)

	for !consumer.IsClosed() {
		event := consumer.Poll(100)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			retryDuration := k.c.Kafka.RetryInitialInterval.AsDuration()

			for attempt := 1; ; attempt++ {
				if err := k.handleMessage(ctx, e); err != nil && !errors.Is(err, ErrInvalidMessage) {
					k.log.WarnContext(ctx,
						"failed to handle Kafka message, retrying",
						slog.Int("attempt", attempt),
						slog.Duration("retry_delay", retryDuration),
						slog.Any("topic_partition", e.TopicPartition),
						slog.Any("error", err),
					)

					time.Sleep(retryDuration)

					retryDuration = min(
						time.Duration(float64(retryDuration)*k.c.Kafka.RetryMultiplier),
						k.c.Kafka.RetryMaxInterval.AsDuration(),
					)

					continue
				}

				break
			}

			consumer.StoreMessage(e)

		case kafka.Error:
			if e.IsFatal() {
				k.log.ErrorContext(ctx,
					"fatal Kafka consumer error",
					slog.Any("error", e),
				)
				return errors.New(e.Error())
			}

			k.log.WarnContext(ctx,
				"Kafka consumer error",
				slog.Any("error", e),
			)
		}
	}

	return nil
}

func (k *KafkaConsumer) Stop(ctx context.Context) error {
	if k.consumer != nil {
		if err := k.consumer.Close(); err != nil {
			k.log.ErrorContext(ctx,
				"failed to close Kafka consumer",
				slog.Any("error", err),
			)
			return err
		}

		k.log.InfoContext(ctx,
			"Kafka consumer stopped",
			slog.String("group_id", groupID),
		)
	}

	return nil
}
