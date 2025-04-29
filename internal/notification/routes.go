package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/notification/handler"
	notificationRepo "github.com/mavrk-mose/pay/internal/notification/repository"
	userRepo "github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/middleware"
)

func NewApiHandler(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	userRepository := userRepo.NewUserRepository(db)
	notificationRepository := notificationRepo.NewNotificationRepo(db)

	notificationHandler := handler.NewNotificationHandler(cfg, userRepository, notificationRepository)

	api := r.Group("/api/v1", middleware.AuthMiddleware())
	{
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.POST("/notifications/:id/read", notificationHandler.MarkAsRead)
	}
}
