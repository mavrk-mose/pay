package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(r *gin.Engine, service NotificationService) {
	h := &NotificationHandler{service: service}

	api := r.Group("/api/notifications")
	{
		api.GET("/", h.GetNotifications)
		api.POST("/:id/read", h.MarkAsRead)
	}
}

// GET /api/notifications?page=1&limit=10
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id") // Assuming you set user_id in middleware after auth
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, err := h.service.GetNotifications(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// POST /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	idParam := c.Param("id")

	notificationID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	err = h.service.MarkNotificationAsRead(userID, notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}
