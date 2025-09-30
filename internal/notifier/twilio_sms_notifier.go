package notifier

import (
	"fmt"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMSNotifier struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioSMSNotifier(accountSID, authToken, fromNumber string) *TwilioSMSNotifier {
	return &TwilioSMSNotifier{
		client:     twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSID, Password: authToken}),
		fromNumber: fromNumber,
	}
}

func (t *TwilioSMSNotifier) SendSMS(req SendSMSRequest) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(req.To)
	params.SetFrom(t.fromNumber)
	params.SetBody(req.Message)

	_, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}
