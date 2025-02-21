package service

import (
	"fmt"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSNotifier sends SMS notifications via Twilio
type SMSNotifier struct {
	client *twilio.RestClient
	from   string // Twilio phone number
}

func NewSMSNotifier(accountSID, authToken, from string) *SMSNotifier {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &SMSNotifier{client: client, from: from}
}

func (n *SMSNotifier) Send(userID, title, message string) error {
	// Fetch the user's phone number from the database (mock implementation)
	to := "+1234567890" // Replace with the user's phone number

	// Send the SMS
	params := &api.CreateMessageParams{}
	params.SetFrom(n.from)
	params.SetTo(to)
	params.SetBody(fmt.Sprintf("%s: %s", title, message))

	_, err := n.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	return nil
}
