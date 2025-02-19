package service

import (
	"fmt"
)

// Dispatcher sends notifications through the appropriate channel
type Dispatcher struct {
	notifiers map[string]Notifier // Map of notification channels
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		notifiers: map[string]Notifier{
			"push":  NewPushNotifier(),
			"sms":   NewSMSNotifier("accountSID", "authToken", "+1234567890"),
			"email": NewEmailNotifier("noreply@example.com", "smtp.example.com", "587", "user", "pass"),
		},
	}
}

// SendNotification sends a notification through the user's preferred channel
func (d *Dispatcher) SendNotification(userID, channel, title, message string) error {
	notifier, exists := d.notifiers[channel]
	if !exists {
		return fmt.Errorf("notification channel %s not supported", channel)
	}

	return notifier.Send(userID, title, message)
}
