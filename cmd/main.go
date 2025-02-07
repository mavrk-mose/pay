package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/pkg/db/postgres"
	"github.com/mavrk-mose/pay/utils/logger"
	"github.com/mavrk-mose/pay/utils"
)

func main() {
	r := gin.Default()

	publicKey, err := middleware.LoadPublicKey("public.pem")
	if err != nil {
		panic(err)
	}


	// Allow 20 requests per second, with a burst of 5
	rl := middleware.NewRateLimiter(rate.Limit(20), 5)

	r.use(rl.RateLimitMiddleware())
	
	// Load configuration
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewApiLogger(cfg)
	
	psqlDB, err := postgres.NewPsqlDB(cfg);
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

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
