package notification

import (
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
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
	return notification, nil
}
