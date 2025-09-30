package notifier

type SendSMSRequest struct {
	To      string
	Message string
}

type SendEmailRequest struct {
	To      string
	Subject string
	Body    string
}
