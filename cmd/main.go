package main

import (
	"github.com/mavrk-mose/pay/internal/middleware"
	"github.com/mavrk-mose/pay/internal/user"
	"golang.org/x/time/rate"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	. "github.com/mavrk-mose/pay/pkg/utils"
)

func main() {
	r := gin.Default()

	_, err := middleware.LoadPublicKey("public.pem")
	if err != nil {
		panic(err)
	}

	// Allow 20 requests per second, with a burst of 5
	rl := middleware.NewRateLimiter(rate.Limit(20), 5)

	r.Use(rl.RateLimitMiddleware())

	// Load configuration
	configPath := GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// Initialize logger
	appLogger := NewApiLogger(cfg)

	// Initialize database
	db, err := config.NewPsqlDB(cfg)
	if err != nil {
		panic("Failed to connect to database!")
		return
	}

	// modules
	user.AuthRoute(db)

	// Use cfg for PORT configuration
	PORT := cfg.Server.Port
	if PORT == "" {
		PORT = "8080" // Fallback to a default port if not specified
	}

	// Start the server
	err = r.Run(":" + PORT)
	if err != nil {
		appLogger.Fatalf("Server failed to start: %v", err)
	}

}
