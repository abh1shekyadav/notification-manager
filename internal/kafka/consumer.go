package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(key, value []byte) error

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(m.Key, m.Value); err != nil {
			// TODO: retry logic / DLQ
		}
	}
}
