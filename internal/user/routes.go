package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	auth "github.com/mavrk-mose/pay/internal/user/handler"
	"github.com/mavrk-mose/pay/pkg/middleware"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/internal/user/service"
)

func AuthRoute(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	auth.InitAuth(cfg)

	userService := service.NewUserService(repository.NewUserRepository(db))

	userHandler := auth.NewUserHandler(userService)

	// Authentication routes
	r.GET("/auth/:provider", auth.BeginAuthHandler)
	r.GET("/auth/:provider/callback", userHandler.AuthCallbackHandler)
	r.GET("/auth/logout/:provider", userHandler.LogoutHandler)

	// Admin routes
	adminRoutes := r.Group("/admin/users", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminRoutes.GET("/", userHandler.ListUsers)
		adminRoutes.POST("/:userID/assign-role", userHandler.AssignRole)
		adminRoutes.POST("/:userID/revoke-role", userHandler.RevokeRole)
		adminRoutes.POST("/:userID/ban", userHandler.BanUser)
		adminRoutes.POST("/:userID/unban", userHandler.UnbanUser)
	}
}
