package notification

import (
	"encoding/json"
	"fmt"

	"github.com/abh1shekyadav/notification-manager/internal/kafka"
)

type NotificationService struct {
	repo     NotificationRepository
	producer *kafka.KafkaProducer
}

func NewNotificationService(repo NotificationRepository, producer *kafka.KafkaProducer) *NotificationService {
	return &NotificationService{repo: repo, producer: producer}
}

func (s *NotificationService) Notify(req NotificationRequest) (*Notification, error) {
	switch req.Type {
	case "email":
		var payload EmailPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, err
		}
	case "sms":
		var payload SMSPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", req.Type)
	}
	notif := NewNotification(req)

	// Save to DB
	if err := s.repo.Save(notif); err != nil {
		return nil, err
	}

	// Publish to Kafka
	if err := s.producer.Publish(notif); err != nil {
		return nil, err
	}

	return notif, nil
}

func (s *NotificationService) FindNotificationByID(notificationID string) (*Notification, error) {
	return s.repo.FindByID(notificationID)
}
