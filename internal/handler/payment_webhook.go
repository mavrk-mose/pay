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

	// - Debits the payer wallet
	// - Credits the payee wallet
	// add the double entries when transaction is successful
	if newStatus == "confirmed" {
		// Create the transaction header (initially pending)
		txn := &models.Transaction{
			ExternalRef: externalRef,
			Status:      models.TransactionPending,
			Details:     nil, // Optionally store extra data as JSON (string)
		}

		// Create ledger entries: one debit and one credit
		debitEntry := models.Transaction{
			WalletID:  payerWalletID,
			EntryType: models.Debit,
			Amount:    amount,
			Currency:  currency,
		}
		creditEntry := models.Transaction{
			WalletID:  payeeWalletID,
			EntryType: models.Credit,
			Amount:    amount,
			Currency:  currency,
		}
		entries := []models.Transaction{debitEntry, creditEntry}

		// Save the transaction and its entries atomically
		if err := s.TransactionRepo.CreateTransactionWithEntries(ctx, txn, entries); err != nil {
			return "", err
		}
	}

	if err := h.TransactionRepo.UpdateTransactionStatus(context.Background(), payload.ExternalRef, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
