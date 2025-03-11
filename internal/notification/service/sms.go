package service

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/internal/user/models"

	"github.com/mavrk-mose/pay/internal/notification/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSNotifier sends SMS notifications via Twilio
type SMSNotifier struct {
	client           *twilio.RestClient
	from             string // Twilio phone number from config
	notificationRepo repository.NotificationRepo
	logger           utils.Logger
}

func NewSMSNotifier(
	accountSID,
	authToken,
	from string,
	notificationRepo repository.NotificationRepo,
	logger utils.Logger,
) *SMSNotifier {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &SMSNotifier{
		client:           client,
		from:             from,
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// Send sends a simple SMS without using a template (for backward compatibility)
func (n *SMSNotifier) Send(ctx context.Context, user models.User, templateID string, details map[string]string) error {
	userID := user.ID.String()

	n.logger.Infof("Sending direct SMS to user %s", userID)

	template, err := n.notificationRepo.GetTemplate(ctx, templateID)
	if err != nil {
		n.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)
	n.logger.Debugf("Processed template message: %s", message)

	params := &api.CreateMessageParams{}
	params.SetFrom(n.from)
	params.SetTo(user.PhoneNumber)
	params.SetBody(fmt.Sprintf("%s: %s", template.Title, message))

	resp, err := n.client.Api.CreateMessage(params)
	if err != nil {
		n.logger.Errorf("Failed to send SMS to %s: %v", user.PhoneNumber, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	n.logger.Infof("Direct SMS sent successfully to user %s (Twilio SID: %s)", userID, *resp.Sid)
	return nil
}
