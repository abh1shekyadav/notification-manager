package notification

import (
	"context"
	"encoding/json"
	"log"

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

func (p *KafkaProducer) Publish(notification *Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(context.Background(), kafka.Message{Value: data})
}

func StartConsumer(brokers []string, topic string, repo NotificationRepository) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "notification-consumer-group",
	})

	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka consumer error: %v", err)
				continue
			}

			var notif Notification
			if err := json.Unmarshal(m.Value, &notif); err != nil {
				log.Printf("Invalid notification message: %v", err)
				continue
			}

			// Process notification based on type
			switch notif.Type {
			case "email":
				var payload EmailPayload
				if err := json.Unmarshal(notif.Payload, &payload); err != nil {
					log.Printf("Invalid email payload: %v", err)
					continue
				}
				log.Printf("Sending email to %s: %s", payload.To, payload.Subject)
			case "sms":
				var payload SMSPayload
				if err := json.Unmarshal(notif.Payload, &payload); err != nil {
					log.Printf("Invalid SMS payload: %v", err)
					continue
				}
				log.Printf("Sending SMS to %s: %s", payload.To, payload.Message)
			default:
				log.Printf("Unknown notification type: %s", notif.Type)
			}

			notif.Status = "SENT"
			if err := repo.UpdateStatus(notif.ID, notif.Status); err != nil {
				log.Printf("Failed to update notification status: %v", err)
			}
		}
	}()
}
