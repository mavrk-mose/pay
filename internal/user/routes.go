package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	auth "github.com/mavrk-mose/pay/internal/user/handler"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/internal/user/service"
)

func AuthRoute(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	auth.InitAuth(cfg)

	userService := service.NewUserService(repository.NewUserRepository(db))

	userHandler := auth.NewUserHandler(userService)

	r.GET("/auth/:provider", auth.BeginAuthHandler)
	r.GET("/auth/:provider/callback", userHandler.AuthCallbackHandler)
	r.GET("/auth/logout/:provider", userHandler.LogoutHandler)
}
