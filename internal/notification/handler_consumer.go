package notification

import (
	"encoding/json"
	"log"
	"time"

	"github.com/abh1shekyadav/notification-manager/internal/kafka"
	"github.com/abh1shekyadav/notification-manager/internal/model"
	"github.com/abh1shekyadav/notification-manager/internal/notifier"
)

const maxRetries = 3

func NewConsumerHandler(
	repo NotificationRepository,
	sms notifier.SMSNotifer,
	email notifier.EmailNotifier,
	dlqProducer *kafka.KafkaProducer,
) func(key, value []byte) error {
	return func(key, value []byte) error {
		var notif model.Notification
		if err := json.Unmarshal(value, &notif); err != nil {
			log.Printf("invalid message: %v", err)
			return err
		}

		log.Printf("Processing notification %s of type %s", notif.ID, notif.Type)
		var err error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = processNotification(notif, sms, email)
			if err == nil {
				_ = repo.UpdateStatus(notif.ID, "SENT")
				return nil
			}
			log.Printf("Attempt %d failed for notification %s: %v", attempt, notif.ID, err)
			time.Sleep(time.Duration(1<<uint(attempt-1)) * time.Second) // 1s, 2s, 4s
		}
		_ = repo.UpdateStatus(notif.ID, "FAILED")
		if dlqProducer != nil {
			if err := dlqProducer.Publish(&notif); err != nil {
				log.Printf("Failed to publish to DLQ for notification %s: %v", notif.ID, err)
			} else {
				log.Printf("Notification %s sent to DLQ", notif.ID)
			}
		}
		return err
	}
}

func processNotification(notif model.Notification, sms notifier.SMSNotifer, email notifier.EmailNotifier) error {
	var err error
	switch notif.Type {
	case "sms":
		var payload model.SMSPayload
		if err = json.Unmarshal(notif.Payload, &payload); err == nil {
			err = sms.SendSMS(notifier.SendSMSRequest{
				To:      payload.To,
				Message: payload.Message,
			})
		}
	case "email":
		var payload model.EmailPayload
		if err = json.Unmarshal(notif.Payload, &payload); err == nil {
			err = email.SendEmail(notifier.SendEmailRequest{
				To:      payload.To,
				Subject: payload.Subject,
				Body:    payload.Body,
			})
		}
	default:
		log.Printf("unsupported notification type: %s", notif.Type)
		return nil
	}
	if err != nil {
		log.Printf("error processing notification %s: %v", notif.ID, err)
		return err
	}
	log.Printf("Successfully processed notification %s", notif.ID)
	return nil
}
