package model

import (
	"encoding/json"
	"time"
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
	UpdatedAt time.Time       `json:"updated_at"`
	Retries   int             `json:"retries"`
	LastError string          `json:"last_error"`
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
