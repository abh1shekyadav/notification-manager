package notification

import (
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo     NotificationRepository
	producer *KafkaProducer
}

func NewNotificationService(repo NotificationRepository, producer *KafkaProducer) *NotificationService {
	return &NotificationService{repo: repo, producer: producer}
}

func (s *NotificationService) Notify(req NotificationRequest) (*Notification, error) {
	notification := &Notification{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(notification); err != nil {
		return nil, err
	}
	if err := s.producer.Publish(notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) FindNotificationById(notificationId string) (*Notification, error) {
	return s.repo.FindById(notificationId)
}
