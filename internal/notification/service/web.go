package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
	. "github.com/mavrk-mose/pay/internal/notification/models"
	"log"
	"sync"
)

// NotificationService handles SSE and notifications.
// It uses a concurrent map for clients where each userID maps to its notification channel.
type NotificationService struct {
	clients   cmap.ConcurrentMap[string, chan Notification]
	templates map[string]NotificationTemplate // Loaded from DB in production.
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		clients:   cmap.New[chan Notification](),
		templates: loadTemplates(), // In production, load these from the DB.
	}
}

// loadTemplates loads preconfigured notification templates
func (s *NotificationService) loadTemplates() map[string]NotificationTemplate {
	return map[string]NotificationTemplate{
		"welcome": {
			ID:      "welcome",
			Title:   "Welcome!",
			Message: "Thank you for joining our platform.",
			Type:    "info",
		},
		"payment_success": {
			ID:      "payment_success",
			Title:   "Payment Successful",
			Message: "Your payment has been processed successfully.",
			Type:    "success",
		},
		"alert": {
			ID:      "alert",
			Title:   "Alert!",
			Message: "An important update requires your attention.",
			Type:    "alert",
		},
	}
}

// SSEHandler handles Server-Sent Events (SSE) connections
func (s *NotificationService) SSEHandler(c *gin.Context) {
	userID := c.Param("userID") // Get the user ID from the request

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Create a channel for this client
	notifyChan := make(chan Notification)
	s.mu.Lock()
	s.clients[userID] = notifyChan
	s.mu.Unlock()

	// Ensure the channel is removed when the client disconnects
	defer func() {
		s.mu.Lock()
		delete(s.clients, userID)
		close(notifyChan)
		s.mu.Unlock()
	}()

	// Listen for client disconnection
	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Client %s disconnected", userID)
			return
		case notification := <-notifyChan:
			// Send the notification as an SSE event
			data, _ := json.Marshal(notification)
			_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			if err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}

// SendNotification sends a notification to a specific user
// SSE provides real-time updates to web clients.
func (s *NotificationService) SendNotification(userID, templateID string, details map[string]string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the preconfigured template
	template, exists := s.templates[templateID]
	if !exists {
		return fmt.Errorf("template %s not found", templateID)
	}

	// Replace placeholders in the template message with details
	message := template.Message
	for _, value := range details {
		message = fmt.Sprintf(message, value)
	}

	// Create the notification
	notification := Notification{
		ID:      uuid.New(),
		UserID:  userID,
		Title:   template.Title,
		Message: message,
		Type:    template.Type,
	}

	// Send the notification to the user's channel
	if clientChan, ok := s.clients[userID]; ok {
		clientChan <- notification
	} else {
		return fmt.Errorf("user %s is not connected", userID)
	}

	return nil
}
