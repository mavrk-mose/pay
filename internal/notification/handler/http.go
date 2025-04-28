package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mavrk-mose/pay/config"
	repository "github.com/mavrk-mose/pay/internal/notification/repository"
	service "github.com/mavrk-mose/pay/internal/notification/service"
	userRepo "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type NotificationHandler struct {
	dispatcher service.Dispatcher
	repo       repository.NotificationRepo
	logger    utils.Logger
}

func NewNotificationHandler(
	r *gin.Engine,
	cfg *config.Config,
	userRepository userRepo.UserRepository,
	notificationRepository repository.NotificationRepo,
	logger utils.Logger,
) *NotificationHandler {
	dispatcher := service.NewDispatcher(
		cfg,
		userRepository,
		notificationRepository,
	)

	return &NotificationHandler{
		dispatcher: *dispatcher,
		repo:       notificationRepository,
		logger:     logger,
	}
}

// GET /api/notifications?page=1&limit=10
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id") // Assuming you set user_id in middleware after auth
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notifications, err := h.repo.FetchNotifications(userID, (page-1)*limit, limit)
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

	err = h.repo.UpdateNotificationAsRead(userID, notificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}
