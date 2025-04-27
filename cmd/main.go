package main

import (
	"github.com/mavrk-mose/pay/internal/payment"
	"github.com/mavrk-mose/pay/internal/user"
	"github.com/mavrk-mose/pay/internal/wallet"
	"github.com/mavrk-mose/pay/pkg/middleware"
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
	}

	// Read the SQL file
	schema, err := os.ReadFile("pkg/general.sql")
	if err != nil {
		log.Fatalln("Failed to read SQL file:", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalln("Failed to execute schema:", err)
	}

	// modules
	user.AuthRoute(r, db, cfg)
	payment.NewApiHandler(r, db, cfg)
	wallet.NewApiHandler(r, db)

	PORT := cfg.Server.Port
	if PORT == "" {
		PORT = "8080"
	}

	err = r.Run(":" + PORT)
	if err != nil {
		appLogger.Fatalf("Server failed to start: %v", err)
	}

}
