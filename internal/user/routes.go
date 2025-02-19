package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	auth "github.com/mavrk-mose/pay/internal/user/handler"
)

func AuthRoute(db *sqlx.DB) {
	// Initialize Google OAuth
	// TODO: read these from the yaml config using viper
	auth.InitGoogleOAuth(
		"your-google-client-id",
		"your-google-client-secret",
		"http://localhost:8080/auth/google/callback",
	)

	// Initialize Gin router
	r := gin.Default()

	// Routes
	r.GET("/auth/google/login", auth.HandleGoogleLogin)
	r.GET("/auth/google/callback", func(c *gin.Context) {
		auth.HandleGoogleCallback(c, db)
	})
}
