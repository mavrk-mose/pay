package executor

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/internal/executor/handler"
	"github.com/mavrk-mose/pay/pkg/nats"
)

func NewApiHandler(r *gin.Engine, db *sqlx.DB) {
	// executorService := service.NewExecutorService(repository.NewExecutorRepo(db))
	executorHandler := handler.NewWebhookHandler(&nats.NatsClient{})

	api := r.Group("/api/v1")

	executor := api.Group("/webhook") //validate with signature
	{
		executor.POST("/stripe", executorHandler.StripeWebhookHandler)
		executor.POST("/paypal", executorHandler.PaypalWebhookHandler)
	}
}