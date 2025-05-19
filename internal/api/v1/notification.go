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
		notification: repository.NewNotificationRepo(db, cfg),
	}
}

// GetNotifications GET /api/notifications?page=1&limit=10
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

// MarkAsRead POST /api/notifications/:id/read
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
