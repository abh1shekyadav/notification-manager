package notification

import (
	"encoding/json"
	"fmt"

	"github.com/abh1shekyadav/notification-manager/internal/kafka"
	"github.com/abh1shekyadav/notification-manager/internal/model"
)

type NotificationService struct {
	repo          NotificationRepository
	smsProducer   *kafka.KafkaProducer
	emailProducer *kafka.KafkaProducer
}

func NewNotificationService(repo NotificationRepository, smsProducer *kafka.KafkaProducer, emailProducer *kafka.KafkaProducer) *NotificationService {
	return &NotificationService{repo: repo, smsProducer: smsProducer, emailProducer: emailProducer}
}

func (s *NotificationService) Notify(req model.NotificationRequest) (*model.Notification, error) {
	var producer *kafka.KafkaProducer
	switch req.Type {
	case "email":
		var payload model.EmailPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, err
		}
		producer = s.emailProducer
	case "sms":
		var payload model.SMSPayload
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			return nil, err
		}
		producer = s.smsProducer
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", req.Type)
	}
	notif := NewNotification(req)

	// Save to DB
	if err := s.repo.Save(notif); err != nil {
		return nil, err
	}

	// Publish to Kafka
	if producer != nil {
		if err := producer.Publish(notif); err != nil {
			_ = s.repo.UpdateStatus(notif.ID, "FAILED")
			return nil, err
		}
	}

	return notif, nil
}

func (s *NotificationService) FindNotificationByID(notificationID string) (*model.Notification, error) {
	return s.repo.FindByID(notificationID)
}
