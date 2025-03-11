package service

import (
	"encoding/json"
	"fmt"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mavrk-mose/pay/pkg/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	. "github.com/mavrk-mose/pay/internal/notification/models"
	"github.com/mavrk-mose/pay/internal/notification/repository"
	"time"
)

// NotificationService handles SSE and notifications.
// It uses a concurrent map for clients where each userID maps to its notification channel.
type WebNotifier struct {
	clients   cmap.ConcurrentMap[string, chan Notification]
	repo      repository.NotificationRepo
	logger    utils.Logger
}

func NewWebNotifier(repo repository.NotificationRepo, logger utils.Logger) *WebNotifier {
	return &WebNotifier{
		clients:   cmap.New[chan Notification](),
		repo: 	   repo,
		logger:    logger,
	}
}

// SSEHandler handles Server-Sent Events (SSE) connections
func (s *WebNotifier) SSEHandler(c *gin.Context) {
	userID := c.Param("userID") 
	s.logger.Infof("SSE connection established for user: %s", userID)

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	notifyChan := make(chan Notification)

	s.clients.Set(userID, notifyChan)
	s.logger.Debugf("Created new channel for user %s", userID)

	// Ensure the channel is removed when the client disconnects.
	defer func() {
		s.clients.Remove(userID)
		close(notifyChan)
		s.logger.Debugf("Removed user %s from active clients map", userID)
	}()

	// Listen for client disconnection
	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("Client %s disconnected", userID)
			return
		case notification := <-notifyChan:
			// Send the notification as an SSE event
			data, err := json.Marshal(notification)
			if err != nil {
				s.logger.Errorf("Failed to marshal notification for user %s: %v", userID, err)
				continue
			}
			
			_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			if err != nil {
				s.logger.Errorf("Failed to write SSE data for user %s: %v", userID, err)
				return
			}

			c.Writer.Flush()
			s.logger.Debugf("Notification sent to user %s: %s", userID, notification.ID)
		}
	}
}

// SendNotification sends a notification to a specific user
// SSE provides real-time updates to web clients.
func (s *WebNotifier) Send(ctx context.Context, userID, templateID string, details map[string]string) error {
	s.logger.Infof("Sending notification to user %s using template %s", userID, templateID)

	template, err := s.repo.GetTemplate(ctx, templateID)
	if err != nil {
		s.logger.Errorf("Failed to get template %s: %v", templateID, err)
		return fmt.Errorf("failed to get template: %w", err)
	}

	message := utils.ReplaceTemplatePlaceholders(template.Message, details)
	s.logger.Debugf("Processed template message: %s", message)

	notification := Notification{
		ID:      uuid.New(),
		UserID:  userID,
		Title:   template.Title,
		Message: message,
		Type:    template.Type,
		Time:    time.Now(),
	}

	if notifyChan, ok := s.clients.Get(userID); ok {
		s.logger.Debugf("User %s is connected, sending notification directly", userID)
		notifyChan <- notification
	} else {
		s.logger.Infof("User %s is not connected, storing notification", userID)
		if err := s.repo.StoreNotification(ctx, notification); err != nil {
			s.logger.Errorf("Failed to store notification for user %s: %v", userID, err)
		}
		return fmt.Errorf("user %s is not connected", userID)
	}

	s.logger.Infof("Notification %s successfully queued for user %s", notification.ID, userID)
	return nil
}
