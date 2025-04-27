package main

import (
	"log"
	"os"

	"github.com/dreson4/graceful/v2"
	"github.com/mavrk-mose/pay/internal/notification"
	"github.com/mavrk-mose/pay/internal/payment"
	"github.com/mavrk-mose/pay/internal/user"
	"github.com/mavrk-mose/pay/internal/wallet"
	"github.com/mavrk-mose/pay/pkg/middleware"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	. "github.com/mavrk-mose/pay/pkg/utils"
)

func main() {
	graceful.Initialize()
	
	r := gin.Default()

	_, err := middleware.LoadPublicKey("public.pem")
	if err != nil {
		panic(err)
	}

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

	db, err := config.NewPsqlDB(cfg)
	if err != nil {
		panic("Failed to connect to database!")
	}

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
	notification.NewApiHandler(r, db)

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
