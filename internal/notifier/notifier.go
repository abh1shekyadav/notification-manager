package notifier

type SMSNotifer interface {
	SendSMS(req SendSMSRequest) error
}

type EmailNotifier interface {
	SendEmail(req SendEmailRequest) error
}
