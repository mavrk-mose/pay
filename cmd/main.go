package main

import (
	"log"
	"os"

	"github.com/dreson4/graceful/v2"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/api"
	"github.com/mavrk-mose/pay/pkg/db"
	. "github.com/mavrk-mose/pay/pkg/utils"
)

func main() {
	graceful.Initialize()

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
	appLogger.InitLogger()

	DB, err := db.NewPsqlDB(cfg)
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.MigrateDB(DB)

	server := api.Init(DB, cfg)

	PORT := cfg.Server.Port
	if PORT == "" {
		PORT = "8080"
	}

	err = server.Run(":" + PORT)
	if err != nil {
		appLogger.Fatalf("Server failed to start: %v", err)
	}

	graceful.Wait()
}
