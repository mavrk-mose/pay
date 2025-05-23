package repository

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/pkg/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/notification/models"
	"github.com/mavrk-mose/pay/internal/notification/service"
)

type notificationRepo struct {
	notifiers map[string]service.Notifier
	logger    utils.Logger
	DB        *sqlx.DB
}

func NewNotificationService(db *sqlx.DB, cfg *config.Config) service.NotificationService {
	notifiers := make(map[string]service.Notifier)

	notifiers["push"] = NewPushNotifier(cfg)
	notifiers["sms"] = NewSMSNotifier(cfg)
	notifiers["email"] = NewEmailNotifier(cfg)
    notifiers["web"] = NewWebNotifier()
	
	return &notificationRepo{
		notifiers: notifiers,
		DB:        db,
	}
}

func (s *notificationRepo) SendNotification(ctx context.Context, user models.User, channel, title string, details map[string]string) error {
	s.logger.Infof("Dispatching notification to user %s via %s channel", user.ID.String(), channel)

	notifier, exists := s.notifiers[channel]
	if !exists {
		s.logger.Errorf("Notification channel %s not supported", channel)
		return fmt.Errorf("notification channel %s not supported", channel)
	}

	if err := notifier.Send(ctx, user, title, details); err != nil {
		s.logger.Errorf("Failed to send notification via %s: %v", channel, err)
		return err
	}

	s.logger.Infof("Successfully sent notification to user %s via %s channel", user.ID.String(), channel)
	return nil
}

func (s *notificationRepo) GetTemplate(ctx context.Context, templateID string) (NotificationTemplate, error) {
	var template NotificationTemplate
	query := `SELECT id, title, subject, message, type, channel, variables, metadata FROM templates WHERE id = $1`
	err := s.DB.GetContext(ctx, &template, query, templateID)
	if err != nil {
		return NotificationTemplate{}, fmt.Errorf("failed to get template: %w", err)
	}
	return template, nil
}

func (s *notificationRepo) StoreNotification(ctx context.Context, notification Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, title, message, type, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, user_id, title, message, type, created_at
	`
	err := s.DB.QueryRowx(
		query,
		notification.ID,
		notification.UserID,
		notification.Title,
		notification.Message,
		notification.Type,
	).StructScan(&notification)
	if err != nil {
		return fmt.Errorf("failed to store notification: %w", err)
	}
	return nil
}

func (s *notificationRepo) FetchNotifications(userID string, page int, limit int) ([]Notification, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, user_id, title, message, type, channel, status, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var notifications []Notification
	err := s.DB.Select(&notifications, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *notificationRepo) UpdateNotificationAsRead(userID string, notificationID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = true
		WHERE id = $1 AND user_id = $2
	`

	_, err := s.DB.Exec(query, notificationID, userID)
	return err
}
