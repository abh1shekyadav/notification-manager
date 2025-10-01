package notification

import (
	"encoding/json"
	"log"

	"github.com/abh1shekyadav/notification-manager/internal/model"
	"github.com/abh1shekyadav/notification-manager/internal/notifier"
)

func NewConsumerHandler(
	repo NotificationRepository,
	sms notifier.SMSNotifer,
	email notifier.EmailNotifier,
) func(key, value []byte) error {
	return func(key, value []byte) error {
		var notif model.Notification
		if err := json.Unmarshal(value, &notif); err != nil {
			log.Printf("invalid message: %v", err)
			return err
		}

		log.Printf("Processing notification %s of type %s", notif.ID, notif.Type)

		switch notif.Type {
		case "sms":
			var payload model.SMSPayload
			if err := json.Unmarshal(notif.Payload, &payload); err != nil {
				_ = repo.UpdateStatus(notif.ID, "FAILED")
				return err
			}
			if err := sms.SendSMS(notifier.SendSMSRequest{
				To:      payload.To,
				Message: payload.Message,
			}); err != nil {
				_ = repo.UpdateStatus(notif.ID, "FAILED")
				return err
			}
			_ = repo.UpdateStatus(notif.ID, "SENT")

		case "email":
			var payload model.EmailPayload
			if err := json.Unmarshal(notif.Payload, &payload); err != nil {
				_ = repo.UpdateStatus(notif.ID, "FAILED")
				return err
			}
			if err := email.SendEmail(notifier.SendEmailRequest{
				To:      payload.To,
				Subject: payload.Subject,
				Body:    payload.Body,
			}); err != nil {
				_ = repo.UpdateStatus(notif.ID, "FAILED")
				return err
			}
			_ = repo.UpdateStatus(notif.ID, "SENT")

		default:
			log.Printf("unsupported notification type: %s", notif.Type)
			_ = repo.UpdateStatus(notif.ID, "FAILED")
		}
		return nil
	}
}
