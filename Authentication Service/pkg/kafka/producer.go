package kafka

import (
	"context"
	"fmt"

	kafkago "github.com/segmentio/kafka-go"
)

type RegisterEventProducer interface {
	PublishRegisterOTP(ctx context.Context, event RegisterOTPEvent) error
	Close() error
}

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafkago.Writer{
			Addr:     kafkago.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafkago.LeastBytes{},
		},
	}
}

func (producer *Producer) PublishRegisterOTP(ctx context.Context, event RegisterOTPEvent) error {
	payload, err := event.Marshal()
	if err != nil {
		return err
	}

	message := kafkago.Message{
		Key:   []byte(event.Email),
		Value: payload,
	}

	if err := producer.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("publish register event: %w", err)
	}

	return nil
}

func (producer *Producer) Close() error {
	if producer == nil || producer.writer == nil {
		return nil
	}
	return producer.writer.Close()
}
