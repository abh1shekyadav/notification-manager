package notification

import "database/sql"

type NotificationRepository interface {
	Save(notification *Notification) error
}

type NotificationRepo struct {
	db *sql.DB
}

func NewNotifcationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Save(notification *Notification) error {
	_, err := r.db.Exec(`
		INSERT INTO notifications (id, user_id, type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		notification.ID, notification.UserID, notification.Type,
		notification.Payload, notification.Status, notification.CreatedAt)
	return err
}
