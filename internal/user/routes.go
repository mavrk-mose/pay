package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	auth "github.com/mavrk-mose/pay/internal/user/handler"
	"github.com/mavrk-mose/pay/internal/user/repository"
)

func AuthRoute(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	auth.InitAuth(cfg)

	// Create repository
	userRepo := repository.NewUserRepository(db)

	// Create handler
	userHandler := auth.NewUserHandler(userRepo)

	// Common auth routes
	r.GET("/auth/:provider", auth.BeginAuthHandler)
	r.GET("/auth/:provider/callback", userHandler.AuthCallbackHandler)
}
