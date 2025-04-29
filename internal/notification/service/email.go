package service

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/config"
	notificationRepo "github.com/mavrk-mose/pay/internal/notification/repository"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/pkg/utils"
	"net/smtp"
)

type EmailNotifier struct {
	from             string
	smtpHost         string
	smtpPort         string
	auth             smtp.Auth
	notificationRepo notificationRepo.NotificationRepo
	logger           utils.Logger
}

func NewEmailNotifier(cfg *config.Config, notificationRepo notificationRepo.NotificationRepo) *EmailNotifier {
	auth := smtp.PlainAuth("", cfg.Server.SMTPUser, cfg.Server.SMTPPassword, cfg.Server.SMTPHost)
	return &EmailNotifier{
		from:             cfg.Server.EmailFrom,
		smtpHost:         cfg.Server.SMTPHost,
		smtpPort:         cfg.Server.SMTPPort,
		auth:             auth,
		notificationRepo: notificationRepo,
	}
}

func (n *EmailNotifier) Send(ctx context.Context, user models.User, templateID string, details map[string]string) error {
	userID := user.ID.String()

	if user.Email == "" {
		n.logger.Warnf("User %s has no email registered", userID)
		return fmt.Errorf("user %s has no email registered", userID)
	}

	to := user.Email

	n.logger.Infof("Sending email to user %s", to)

	template, err := n.notificationRepo.GetTemplate(ctx, templateID)
	if err != nil {
		n.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)

	n.logger.Debugf("Processed template message: %s", message)

	subject := fmt.Sprintf("Subject: %s\n", templateID)
	body := fmt.Sprintf("%s\n\n%s", templateID, message)
	msg := []byte(subject + "\n" + body)

	err = smtp.SendMail(n.smtpHost+":"+n.smtpPort, n.auth, n.from, []string{to}, msg)
	if err != nil {
		n.logger.Errorf("Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
