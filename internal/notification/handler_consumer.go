package notification

import (
	"encoding/json"
	"log"
	"time"

	"github.com/abh1shekyadav/notification-manager/internal/kafka"
	"github.com/abh1shekyadav/notification-manager/internal/metrics"
	"github.com/abh1shekyadav/notification-manager/internal/model"
	"github.com/abh1shekyadav/notification-manager/internal/notifier"
)

const maxRetries = 3
const maxBackoff = 10 * time.Second

func NewConsumerHandler(
	repo NotificationRepository,
	sms notifier.SMSNotifer,
	email notifier.EmailNotifier,
	dlqProducer *kafka.KafkaProducer,
	logger *log.Logger,
) func(key, value []byte) error {

	return func(key, value []byte) error {
		var notif model.Notification
		if err := json.Unmarshal(value, &notif); err != nil {
			logger.Printf("invalid message: %v", err)
			return err
		}

		logger.Printf("Processing notification %s of type %s", notif.ID, notif.Type)
		var err error

		for attempt := 1; attempt <= maxRetries; attempt++ {
			start := time.Now()
			err = processNotification(notif, sms, email, logger)
			duration := time.Since(start).Seconds()

			notifType := notif.Type
			if notifType != "sms" && notifType != "email" {
				notifType = "unknown"
			}

			metrics.ProcessingDuration.WithLabelValues(notifType).Observe(duration)

			if err == nil {
				_ = repo.UpdateStatus(notif.ID, "SENT")
				metrics.NotificationsTotal.WithLabelValues(notifType, "SENT").Inc()
				return nil
			}

			logger.Printf("Attempt %d failed for notification %s: %v", attempt, notif.ID, err)
			sleep := time.Duration(1<<uint(attempt-1)) * time.Second
			if sleep > maxBackoff {
				sleep = maxBackoff
			}
			time.Sleep(sleep)
		}

		_ = repo.UpdateStatus(notif.ID, "FAILED")
		metrics.NotificationsTotal.WithLabelValues(notif.Type, "FAILED").Inc()

		if dlqProducer != nil {
			if err := dlqProducer.Publish(&notif); err != nil {
				logger.Printf("Failed to publish to DLQ for notification %s: %v", notif.ID, err)
			} else {
				logger.Printf("Notification %s sent to DLQ", notif.ID)
			}
		}

		return err
	}
}

func processNotification(notif model.Notification, sms notifier.SMSNotifer, email notifier.EmailNotifier, logger *log.Logger) error {
	var err error

	switch notif.Type {
	case "sms":
		if sms == nil {
			logger.Printf("SMS notifier not configured, skipping")
			return nil
		}
		var payload model.SMSPayload
		if err = json.Unmarshal(notif.Payload, &payload); err == nil {
			err = sms.SendSMS(notifier.SendSMSRequest{
				To:      payload.To,
				Message: payload.Message,
			})
		}
	case "email":
		if email == nil {
			logger.Printf("Email notifier not configured, skipping")
			return nil
		}
		var payload model.EmailPayload
		if err = json.Unmarshal(notif.Payload, &payload); err == nil {
			err = email.SendEmail(notifier.SendEmailRequest{
				To:      payload.To,
				Subject: payload.Subject,
				Body:    payload.Body,
			})
		}
	default:
		logger.Printf("unsupported notification type: %s", notif.Type)
		return nil
	}

	if err != nil {
		logger.Printf("error processing notification %s: %v", notif.ID, err)
		return err
	}

	logger.Printf("Successfully processed notification %s", notif.ID)
	return nil
}
