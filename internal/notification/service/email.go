package service

import (
	"context"
	"fmt"
	notificationRepo "github.com/mavrk-mose/pay/internal/notification/repository"
	userRepo "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
	"net/smtp"
)

// EmailNotifier sends email notifications via SMTP
type EmailNotifier struct {
	from             string // Sender email address
	smtpHost         string // SMTP server host
	smtpPort         string // SMTP server port
	auth             smtp.Auth
	userRepo         userRepo.UserRepository
	notificationRepo notificationRepo.NotificationRepo
	logger           utils.Logger
}

func NewEmailNotifier(
	from,
	smtpHost,
	smtpPort,
	smtpUser,
	smtpPass string,
	userRepo userRepo.UserRepository,
	notificationRepo notificationRepo.NotificationRepo,
	logger utils.Logger,
) *EmailNotifier {
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return &EmailNotifier{
		from:             from,
		smtpHost:         smtpHost,
		smtpPort:         smtpPort,
		auth:             auth,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

func (n *EmailNotifier) Send(ctx context.Context, userID, templateID string, details map[string]string) error {
	user, err := n.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		n.logger.Errorf("Failed to get user %s: %v", userID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.Email == "" {
		n.logger.Warnf("User %s has no email registered", userID)
		return fmt.Errorf("user %s has no email registered", userID)
	}

	to := user.Email

	template, err := n.notificationRepo.GetTemplate(ctx, templateID)
	if err != nil {
		n.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)

	subject := fmt.Sprintf("Subject: %s\n", templateID)
	body := fmt.Sprintf("%s\n\n%s", templateID, message)
	msg := []byte(subject + "\n" + body)

	// Send the email
	err = smtp.SendMail(n.smtpHost+":"+n.smtpPort, n.auth, n.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
