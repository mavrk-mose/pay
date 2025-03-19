package handler

import (
	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	repository "github.com/mavrk-mose/pay/internal/ledger/repository"
	ledgerService "github.com/mavrk-mose/pay/internal/ledger/service"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	paymentService "github.com/mavrk-mose/pay/internal/payment/service"
	"net/http"
	"time"
)

type PaymentHandler struct {
	paymentService paymentService.PaymentService
	ledgerService  ledgerService.LedgerService
	txnRepo        repository.Repo
}

type WebhookPayload struct {
	ExternalRef string `json:"external_ref"`
	Status      string `json:"status"` // e.g., "confirmed" or "failed"
}

func NewPaymentHandler(
	paymentService paymentService.PaymentService,
	ledgerService ledgerService.LedgerService,
	txnRepo repository.Repo,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		ledgerService:  ledgerService,
		txnRepo:        txnRepo,
	}
}

func (t *PaymentHandler) Check(c *gin.Context) {
	c.Status(200)
}

func (t *PaymentHandler) ProcessPayment(ctx *gin.Context) {
	var paymentIntent PaymentIntent

	if err := ctx.ShouldBindJSON(&paymentIntent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := t.paymentService.ProcessPayment(ctx, paymentIntent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}

func (t *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")

	status, err := t.paymentService.GetPaymentStatus(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// QueryPayments allows filtering payments based on type, date range, and status
func (t *PaymentHandler) QueryPayments(c *gin.Context) {
	userID := c.Param("userID")
	status := c.Query("status")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	var startDate, endDate time.Time
	var err error

	// Parse date range if provided
	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use RFC3339"})
			return
		}

		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use RFC3339"})
			return
		}
	}

	// Parse payment status if provided
	var paymentStatus PaymentStatus
	if status != "" {
		if err := paymentStatus.UnmarshalText([]byte(status)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
			return
		}
	}

	// Query payments based on filters
	var payments []PaymentIntent

	// Filter by date range if applicable
	if !startDate.IsZero() && !endDate.IsZero() {
		payments, err = t.paymentService.QueryPaymentsByDateRange(userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Filter by status if applicable
	if status != "" {
		payments, err = t.paymentService.QueryPaymentsByStatus(userID, paymentStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return filtered payments
	c.JSON(http.StatusOK, payments)
}

func (t *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")
	var request struct {
		Status PaymentStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := t.paymentService.UpdatePaymentStatus(paymentID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}

func (t *PaymentHandler) GetPaymentDetails(c *gin.Context) {
	paymentID := c.Param("paymentID")

	payment, err := t.paymentService.GetPaymentDetails(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var payload WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	txn, err := h.ledgerService.GetTransactionByID(payload.ExternalRef)
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
			err := h.ledgerService.RecordTransaction(c, txn)
			if err != nil {

			}
		}()
	}

	if err := h.txnRepo.UpdateTransactionStatus(c, payload.ExternalRef, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
