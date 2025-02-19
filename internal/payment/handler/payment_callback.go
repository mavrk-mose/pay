package handler

import (
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	repository "github.com/mavrk-mose/pay/internal/ledger/repository"
	"github.com/mavrk-mose/pay/internal/ledger/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentWebhookHandler struct {
	TransactionRepo repository.Repo
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

	txn, err := h.ledger.GetTransactionByID(payload.ExternalRef)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var newStatus TransactionStatus
	switch txn.Status {
	case "confirmed":
		newStatus = TransactionConfirmed
	case "failed":
		newStatus = TransactionFailed
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown status"})
		return
	}

	if newStatus == "confirmed" {
		go func() {
			err := h.ledger.RecordTransaction(c, txn)
			if err != nil {

			}
		}()
	}

	if err := h.TransactionRepo.UpdateTransactionStatus(c, payload.ExternalRef, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
