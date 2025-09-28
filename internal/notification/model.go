package notification

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationRequest struct {
	UserID  string          `json:"user_id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Notification struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewNotification(req NotificationRequest) *Notification {
	return &Notification{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Type:      req.Type,
		Payload:   req.Payload,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}
}

type EmailPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

type SMSPayload struct {
	Message string `json:"message"`
	To      string `json:"to"`
}
