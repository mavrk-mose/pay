package service

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// PushNotifier sends push notifications via Firebase Cloud Messaging (FCM)
type PushNotifier struct{}

func NewPushNotifier() *PushNotifier {
	return &PushNotifier{}
}

func (n *PushNotifier) Send(userID, title, message string) error {
	// Initialize Firebase app
	ctx := context.Background()
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	// Get the messaging client
	client, err := app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("failed to get messaging client: %v", err)
	}

	// Send the notification
	_, err = client.Send(ctx, &messaging.Message{
		Token: "user_device_token", // Replace with the user's device token
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send push notification: %v", err)
	}

	return nil
}
