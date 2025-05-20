package repository

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/mavrk-mose/pay/internal/notification/models"
	"github.com/mavrk-mose/pay/internal/notification/service"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/pkg/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

// WebNotifier handles SSE and notifications.
// It uses a concurrent map for clients where each userID maps to its notification channel.
type WebNotifier struct {
	clients      cmap.ConcurrentMap[string, chan Notification]
	notification service.NotificationService
	logger       utils.Logger
}

func NewWebNotifier() *WebNotifier {
	return &WebNotifier{
		clients: cmap.New[chan Notification](),
	}
}

// SSEHandler handles Server-Sent Events (SSE) connections
func (s *WebNotifier) SSEHandler(c *gin.Context) {
	userID := c.Param("userID")
	s.logger.Infof("SSE connection established for user: %s", userID)

	_, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("X-Accel-Buffering: no")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	notifyChan := make(chan Notification)

	s.clients.Set(userID, notifyChan)
	s.logger.Debugf("Created new channel for user %s", userID)

	defer func() {
		s.clients.Remove(userID)
		close(notifyChan)
		s.logger.Debugf("Removed user %s from active clients map", userID)
	}()

	ctx := c
	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("Client %s disconnected", userID)
			return
		case notification := <-notifyChan:
			c.SSEvent("message", notification.Message)
			c.Writer.Flush()
			s.logger.Debugf("Notification sent to user %s: %s", userID, notification.ID)
		}
	}
}

// Send sends a notification to a specific user
// SSE provides real-time updates to web clients.
func (s *WebNotifier) Send(ctx context.Context, user models.User, templateID string, details map[string]string) error {
	userID := user.ID.String()

	s.logger.Infof("Sending notification to user %s using template %s", userID, templateID)

	template, err := s.notification.GetTemplate(ctx, templateID)
	if err != nil {
		s.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)

	s.logger.Debugf("Processed template message: %s", message)

	notification := Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Channel:   "WEB",
		Title:     template.Title,
		Message:   message,
		Type:      template.Type,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if notifyChan, ok := s.clients.Get(userID); ok {
		s.logger.Debugf("User %s is connected, sending notification directly", userID)
		notifyChan <- notification
	} else {
		s.logger.Infof("User %s is not connected, storing notification", userID)
		if err := s.notification.StoreNotification(ctx, notification); err != nil {
			s.logger.Errorf("Failed to store notification for user %s: %v", userID, err)
		}
		return fmt.Errorf("user %s is not connected", userID)
	}

	s.logger.Infof("Notification %s successfully queued for user %s", notification.ID, userID)
	return nil
}
