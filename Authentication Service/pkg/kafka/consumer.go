package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	kafkago "github.com/segmentio/kafka-go"
)

type RegisterEventHandler func(ctx context.Context, event RegisterOTPEvent) error

type RegisterEventConsumer struct {
	reader  *kafkago.Reader
	handler RegisterEventHandler
}

func NewRegisterEventConsumer(brokers []string, topic string, groupID string, handler RegisterEventHandler) *RegisterEventConsumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &RegisterEventConsumer{
		reader:  reader,
		handler: handler,
	}
}

func (consumer *RegisterEventConsumer) Start(ctx context.Context) {
	for {
		message, err := consumer.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Error("kafka consumer fetch error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		event, err := ParseRegisterOTPEvent(message.Value)
		if err != nil {
			logger.Error("kafka register event parse error: %v", err)
			if commitErr := consumer.reader.CommitMessages(ctx, message); commitErr != nil {
				logger.Error("kafka commit parse-failed message error: %v", commitErr)
			}
			continue
		}

		if err := consumer.handler(ctx, event); err != nil {
			logger.Error("kafka register event handler error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err := consumer.reader.CommitMessages(ctx, message); err != nil {
			logger.Error("kafka commit register event error: %v", err)
			time.Sleep(2 * time.Second)
		}
	}
}

func (consumer *RegisterEventConsumer) Close() error {
	if consumer == nil || consumer.reader == nil {
		return nil
	}
	if err := consumer.reader.Close(); err != nil {
		return fmt.Errorf("close kafka reader: %w", err)
	}
	return nil
}
