package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/pkg/middleware"
	"github.com/mavrk-mose/pay/config"

	handler "github.com/mavrk-mose/pay/internal/notification/handler"
	userRepo "github.com/mavrk-mose/pay/internal/user/repository"
	notificationRepo "github.com/mavrk-mose/pay/internal/notification/repository"
)

func NewApiHandler(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	userRepository := userRepo.NewUserRepository(db)                
	notificationRepository := notificationRepo.NewNotificationRepo(db) 

	notificationHandler := handler.NewNotificationHandler(r, cfg, userRepository, notificationRepository)

	// Notification Routes
	api := r.Group("/api/v1", middleware.AuthMiddleware())
	{
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.POST("/notifications/:id/read", notificationHandler.MarkAsRead)
	}
}
