package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/notification/models"
)

type NotificationRepo interface {
	GetTemplate(ctx context.Context, templateID string) (NotificationTemplate, error)
	StoreNotification(ctx context.Context, notification Notification) error
	FetchNotifications(userID string, i int, limit int) ([]Notification, error)
	UpdateNotificationAsRead(userID string, notificationID uuid.UUID) error
}

type notificationRepo struct {
	DB *sqlx.DB
}

func NewNotificationRepo(db *sqlx.DB) NotificationRepo {
	return &notificationRepo{DB: db}
}

func (s *notificationRepo) GetTemplate(ctx context.Context, templateID string) (NotificationTemplate, error) {
	var template NotificationTemplate
	query := `SELECT id, title, message, type FROM templates WHERE id = $1`
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

