package notification

import "time"

type NotificationRequest struct {
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
