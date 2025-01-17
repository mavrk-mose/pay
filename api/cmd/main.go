package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/common/config"
	"github.com/mavrk-mose/pay/api/pkg/db/postgres"
	"github.com/mavrk-mose/pay/common/logger"
	"github.com/mavrk-mose/pay/common/utils"
)

func main() {
	r := gin.Default()
	
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
