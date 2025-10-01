package notification

import (
	"database/sql"
	"encoding/json"

	"github.com/abh1shekyadav/notification-manager/internal/model"
)

type NotificationRepository interface {
	Save(notification *model.Notification) error
	FindByID(notificationID string) (*model.Notification, error)
	UpdateStatus(notificationID string, status string) error
}

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Save(notification *model.Notification) error {
	_, err := r.db.Exec(`
		INSERT INTO notifications (id, user_id, type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		notification.ID, notification.UserID, notification.Type,
		string(notification.Payload), notification.Status, notification.CreatedAt)
	return err
}

func (r *NotificationRepo) FindByID(notificationID string) (*model.Notification, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, type, payload, status, created_at
		FROM notifications WHERE id = $1`, notificationID)

	var notification model.Notification
	var payloadText string
	err := row.Scan(&notification.ID, &notification.UserID, &notification.Type,
		&payloadText, &notification.Status, &notification.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	notification.Payload = json.RawMessage(payloadText)
	return &notification, nil
}

func (r *NotificationRepo) UpdateStatus(notificationID string, status string) error {
	_, err := r.db.Exec(`
		UPDATE notifications SET status = $1, updated_at = now() WHERE id = $2`, status, notificationID)
	return err
}
