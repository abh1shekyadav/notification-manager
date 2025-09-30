package kafka

import (
	"context"
	"encoding/json"

	"github.com/abh1shekyadav/notification-manager/internal/notification"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Publish(notification *notification.Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
}
