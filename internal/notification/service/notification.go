package service

import (
	"context"
	"github.com/google/uuid"

	. "github.com/mavrk-mose/pay/internal/notification/models"
	"github.com/mavrk-mose/pay/internal/user/models"
)

//go:generate mockery --name=NotificationService --output=./mocks --filename=notification.go --with-expecter
type NotificationService interface {
	SendNotification(ctx context.Context, user models.User, channel, title string, details map[string]string) error
	GetTemplate(ctx context.Context, templateID string) (NotificationTemplate, error)
	StoreNotification(ctx context.Context, notification Notification) error
	FetchNotifications(userID string, i int, limit int) ([]Notification, error)
	UpdateNotificationAsRead(userID string, notificationID uuid.UUID) error
}
