package main

import (
	"log"
	"os"

	"github.com/dreson4/graceful/v2"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/api/middleware"
	"github.com/mavrk-mose/pay/internal/api"
	"github.com/mavrk-mose/pay/pkg/db"
	. "github.com/mavrk-mose/pay/pkg/utils"
)

func main() {
	graceful.Initialize()

	r := gin.Default()

	rl := middleware.NewRateLimiter(rate.Limit(20), 5)
	r.Use(rl.RateLimitMiddleware())

	configPath := GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := NewApiLogger(cfg)

	DB, err := db.NewPsqlDB(cfg)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.MigrateDB(DB)

	// modules
	api.NewApiHandler(r, DB, cfg)

	PORT := cfg.Server.Port
	if PORT == "" {
		PORT = "8080"
	}

	err = r.Run(":" + PORT)
	if err != nil {
		appLogger.Fatalf("Server failed to start: %v", err)
	}

	graceful.Wait()
}
