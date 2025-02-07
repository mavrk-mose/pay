package handler

import (
	"context"
	"net/http"
	"github.com/mavrk-mose/pay/internal/models"
	"github.com/mavrk-mose/pay/internal/repository"

	"github.com/gin-gonic/gin"
)

type PaymentWebhookHandler struct {
	TransactionRepo repository.TransactionRepo
	ledger          service.LedgerService
}

type WebhookPayload struct {
	ExternalRef string `json:"external_ref"`
	Status      string `json:"status"` // e.g., "confirmed" or "failed"
}

func (h *PaymentWebhookHandler) HandleWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	var newStatus models.TransactionStatus
	switch payload.Status {
	case "confirmed":
		newStatus = models.TransactionConfirmed
	case "failed":
		newStatus = models.TransactionFailed
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status"})
		return
	}

	if newStatus == "confirmed" {
		go ledger.RecordTransaction(payload)
	}

	if err := h.TransactionRepo.UpdateTransactionStatus(context.Background(), payload.ExternalRef, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
