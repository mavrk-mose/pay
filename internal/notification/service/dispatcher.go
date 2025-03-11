package service

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/notification/repository"
	user "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type Dispatcher struct {
	notifiers map[string]Notifier
	logger    utils.Logger
}

func NewDispatcher(
	cfg *config.Config,
	userRepo user.UserRepository,
	notificationRepo repository.NotificationRepo,
	logger utils.Logger,
) *Dispatcher {
	logger.Info("Initializing notification dispatcher")

	notifiers := make(map[string]Notifier)

	notifiers["push"] = NewPushNotifier(
		notificationRepo,
		userRepo,
		logger,
		cfg.Firebase,
	)
	logger.Debug("Push notifier initialized")

	if cfg.Twilio.AccountSID != "" && cfg.Twilio.AuthToken != "" && cfg.Twilio.From != "" {

		notifiers["sms"] = NewSMSNotifier(
			cfg.Twilio.AccountSID,
			cfg.Twilio.AuthToken,
			cfg.Twilio.From,
			userRepo,
			notificationRepo,
			logger,
		)

		logger.Infof("SMS notifier initialized with Twilio from: %s", cfg.Twilio.From)
	} else {
		logger.Warn("SMS notifier not initialized: missing Twilio configuration")
	}

	if cfg.Server.EmailFrom != "" && cfg.Server.SMTPHost != "" {
		notifiers["email"] = NewEmailNotifier(
			cfg.Server.EmailFrom,
			cfg.Server.SMTPHost,
			cfg.Server.SMTPPort,
			cfg.Server.SMTPUser,
			cfg.Server.SMTPPassword,
			userRepo,
			notificationRepo,
			logger,
		)
		logger.Infof("Email notifier initialized with sender: %s", cfg.Server.EmailFrom)
	} else {
		logger.Warn("Email notifier not initialized: missing email configuration")
	}

	// Add Web Notifier if configured
	// Assuming webhook URL is in the config
	if cfg.Server.WebhookURL != "" {
		notifiers["web"] = NewWebNotifier(notificationRepo, logger)
		logger.Infof("Web notifier initialized with webhook URL: %s", cfg.Server.WebhookURL)
	} else {
		logger.Warn("Web notifier not initialized: missing webhook configuration")
	}

	return &Dispatcher{
		notifiers: notifiers,
		logger:    logger,
	}
}

// SendNotification Dispatcher sends notifications through the appropriate channel
// SendNotification sends a notification through the user's preferred channel
func (d *Dispatcher) SendNotification(ctx context.Context, userID, channel, title string, details map[string]string) error {
	d.logger.Infof("Dispatching notification to user %s via %s channel", userID, channel)

	notifier, exists := d.notifiers[channel]
	if !exists {
		d.logger.Errorf("Notification channel %s not supported", channel)
		return fmt.Errorf("notification channel %s not supported", channel)
	}

	if err := notifier.Send(ctx, userID, title, details); err != nil {
		d.logger.Errorf("Failed to send notification via %s: %v", channel, err)
		return err
	}

	d.logger.Infof("Successfully sent notification to user %s via %s channel", userID, channel)
	return nil
}
