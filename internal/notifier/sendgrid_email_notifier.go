package notifier

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridEmailNotifier struct {
	client    *sendgrid.Client
	fromEmail string
}

func NewSendGridEmailNotifier(apiKey, fromEmail string) *SendGridEmailNotifier {
	return &SendGridEmailNotifier{
		client:    sendgrid.NewSendClient(apiKey),
		fromEmail: fromEmail,
	}
}

func (s *SendGridEmailNotifier) SendEmail(req SendEmailRequest) error {
	from := mail.NewEmail("Notification Service", s.fromEmail)
	to := mail.NewEmail("", req.To)
	message := mail.NewSingleEmail(from, req.Subject, to, req.Body, req.Body)
	_, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
