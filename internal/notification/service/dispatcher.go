package service

import (
	"context"
	"fmt"

	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/notification/repository"
	"github.com/mavrk-mose/pay/internal/user/models"
	user "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type DispatcherService interface {
	SendNotification(ctx context.Context, user models.User, channel, title string, details map[string]string) error
}

type Dispatcher struct {
	notifiers map[string]Notifier
	logger    utils.Logger
}

func NewDispatcher(
	cfg *config.Config,
	userRepo user.UserRepository,
	notificationRepo repository.NotificationRepo,
) *Dispatcher {
	notifiers := make(map[string]Notifier)

	notifiers["push"] = NewPushNotifier(
		notificationRepo,
		userRepo,
		cfg,
	)

	if cfg.Twilio.AccountSID != "" && cfg.Twilio.AuthToken != "" && cfg.Twilio.From != "" {
		notifiers["sms"] = NewSMSNotifier(cfg, notificationRepo)
	}

	if cfg.Server.EmailFrom != "" && cfg.Server.SMTPHost != "" {
		notifiers["email"] = NewEmailNotifier(cfg, notificationRepo)
	}

	if cfg.Server.WebhookURL != "" {
		notifiers["web"] = NewWebNotifier(notificationRepo)
	}

	return &Dispatcher{
		notifiers: notifiers,
	}
}

func (d *Dispatcher) SendNotification(ctx context.Context, user models.User, channel, title string, details map[string]string) error {
	d.logger.Infof("Dispatching notification to user %s via %s channel", user.ID.String(), channel)

	notifier, exists := d.notifiers[channel]
	if !exists {
		d.logger.Errorf("Notification channel %s not supported", channel)
		return fmt.Errorf("notification channel %s not supported", channel)
	}

	if err := notifier.Send(ctx, user, title, details); err != nil {
		d.logger.Errorf("Failed to send notification via %s: %v", channel, err)
		return err
	}

	d.logger.Infof("Successfully sent notification to user %s via %s channel", user.ID.String(), channel)
	return nil
}
