package service

import (
	"context"
	"fmt"

	"github.com/mavrk-mose/pay/internal/notification/repository"
	repo "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSNotifier sends SMS notifications via Twilio
type SMSNotifier struct {
	client           *twilio.RestClient
	from             string // Twilio phone number from config
	userRepo         repo.UserRepository
	notificationRepo repository.NotificationRepo
	logger           utils.Logger
}

func NewSMSNotifier(
	accountSID,
	authToken,
	from string,
	userRepo repo.UserRepository,
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
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// SendSMS sends an SMS notification using a template
func (n *SMSNotifier) SendSMS(ctx context.Context, userID, title string, details map[string]string) error {
	n.logger.Infof("Preparing SMS for user %s using template %s", userID, title)

	user, err := n.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		n.logger.Errorf("Failed to get user %s: %v", userID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.PhoneNumber == "" {
		n.logger.Warnf("User %s has no phone number registered", userID)
		return fmt.Errorf("user %s has no phone number registered", userID)
	}

	template, err := n.notificationRepo.GetTemplate(ctx, title)
	if err != nil {
		n.logger.Errorf("Failed to get template %s: %v", title, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)

	params := &api.CreateMessageParams{}
	params.SetFrom(n.from)
	params.SetTo(user.PhoneNumber)
	params.SetBody(fmt.Sprintf("%s: %s", template.Title, message))

	n.logger.Debugf("Sending SMS to %s with message: %s", user.PhoneNumber, message)

	resp, err := n.client.Api.CreateMessage(params)
	if err != nil {
		n.logger.Errorf("Failed to send SMS to %s: %v", user.PhoneNumber, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	n.logger.Infof("SMS sent successfully to user %s (Twilio SID: %s)", userID, *resp.Sid)
	return nil
}

// Send sends a simple SMS without using a template (for backward compatibility)
func (n *SMSNotifier) Send(ctx context.Context, userID, templateID string, details map[string]string) error {
	n.logger.Infof("Sending direct SMS to user %s", userID)

	user, err := n.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		n.logger.Errorf("Failed to get user %s: %v", userID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.PhoneNumber == "" {
		n.logger.Warnf("User %s has no phone number registered", userID)
		return fmt.Errorf("user %s has no phone number registered", userID)
	}

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

	// Send the SMS
	resp, err := n.client.Api.CreateMessage(params)
	if err != nil {
		n.logger.Errorf("Failed to send SMS to %s: %v", user.PhoneNumber, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	n.logger.Infof("Direct SMS sent successfully to user %s (Twilio SID: %s)", userID, *resp.Sid)
	return nil
}
