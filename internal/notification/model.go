package notification

import (
	"time"

	"github.com/abh1shekyadav/notification-manager/internal/model"
	"github.com/google/uuid"
)

func NewNotification(req model.NotificationRequest) *model.Notification {
	return &model.Notification{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}
}
