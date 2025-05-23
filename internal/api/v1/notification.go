package handler

import (
	"github.com/mavrk-mose/pay/internal/notification/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/notification/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type NotificationHandler struct {
	notification service.NotificationService
	logger       utils.Logger
}

func NewNotificationHandler(cfg *config.Config, db *sqlx.DB) *NotificationHandler {
	return &NotificationHandler{
		notification: repository.NewNotificationService(db, cfg),
	}
}

// GetNotifications godoc
// @Summary      Get notifications
// @Description  Retrieves notifications for a user
// @Tags         notifications
// @Produce      json
// @Param        page   query  int  false  "Page number"
// @Param        limit  query  int  false  "Page size"
// @Success      200  {array}  models.Notification
// @Failure      500  {object}  map[string]string
// @Router       /api/notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	h.logger.Infof("Fetching notifications for user %s page %d limit %d", userID, page, limit)

	notifications, err := h.notification.FetchNotifications(userID, (page-1)*limit, limit)
	if err != nil {
		h.logger.Errorf("Failed to fetch notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}

	h.logger.Infof("Fetched %d notifications for user %s", len(notifications), userID)

	c.JSON(http.StatusOK, notifications)
}

// MarkAsRead godoc
// @Summary      Mark notification as read
// @Description  Marks a notification as read for a user
// @Tags         notifications
// @Produce      json
// @Param        id  path  string  true  "Notification ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	idParam := c.Param("id")

	h.logger.Infof("Marking notification %s as read for user %s", idParam, userID)

	notificationID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Errorf("Invalid notification ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	err = h.notification.UpdateNotificationAsRead(userID, notificationID)
	if err != nil {
		h.logger.Errorf("Failed to mark notification as read: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	h.logger.Infof("Notification %s marked as read for user %s", idParam, userID)

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

// SSEHandler godoc
// @Summary      Notification SSE stream
// @Description  Establishes a Server-Sent Events (SSE) connection for real-time notifications
// @Tags         notifications
// @Produce      text/event-stream
// @Success      200  {string}  string  "SSE stream established"
// @Failure      401  {object}  map[string]string
// @Router       /api/notifications/stream [get]
func (h *NotificationHandler) SSEHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id required in context"})
		return
	}

	h.logger.Infof("Establishing SSE connection for user %s", userID)

	type webNotifierProvider interface {
		WebNotifier() *repository.WebNotifier
	}

	provider, ok := h.notification.(webNotifierProvider)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "SSE not supported"})
		return
	}
	webNotifier := provider.WebNotifier()

	notifyChan := webNotifier.RegisterClient(userID)
	defer webNotifier.UnregisterClient(userID)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	ctx := c.Request.Context()
	for {
		select {
		case notification := <-notifyChan:
			c.SSEvent("notification", notification)
			c.Writer.Flush()
		case <-ctx.Done():
			return
		}
	}
}
