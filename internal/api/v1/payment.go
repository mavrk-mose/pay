package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	repository "github.com/mavrk-mose/pay/internal/ledger/repository"
	ledger "github.com/mavrk-mose/pay/internal/ledger/service"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	payment "github.com/mavrk-mose/pay/internal/payment/service"
)

type PaymentHandler struct {
	db             *sqlx.DB
	paymentService payment.PaymentService
	ledgerService  ledger.LedgerService
}

type WebhookPayload struct {
	ExternalRef string `json:"external_ref"`
	Status      string `json:"status"` // e.g., "confirmed" or "failed"
}

func NewPaymentHandler(db *sqlx.DB) *PaymentHandler {
	return &PaymentHandler{
		ledgerService: repository.NewLedgerService(db),
	}
}

func (h *PaymentHandler) Check(c *gin.Context) {
	c.Status(200)
}

// ProcessPayment godoc
// @Summary      Process a payment
// @Description  Processes a payment intent and returns a success message
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        paymentIntent  body  PaymentIntent  true  "Payment Intent"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/payments [post]
func (h *PaymentHandler) ProcessPayment(ctx *gin.Context) {
	var paymentIntent PaymentIntent

	if err := ctx.ShouldBindJSON(&paymentIntent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentService.ProcessPayment(ctx, paymentIntent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}

// GetPaymentStatus godoc
// @Summary      Get payment status
// @Description  Returns the status of a payment by ID
// @Tags         payments
// @Produce      json
// @Param        paymentID  path  string  true  "Payment ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/payments/{paymentID}/status [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")

	status, err := h.paymentService.GetPaymentStatus(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

// QueryPayments allows filtering payments based on type, date range, and status
func (h *PaymentHandler) QueryPayments(c *gin.Context) {
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
		payments, err = h.paymentService.QueryPaymentsByDateRange(userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Filter by status if applicable
	if status != "" {
		payments, err = h.paymentService.QueryPaymentsByStatus(userID, paymentStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return filtered payments
	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")
	var request struct {
		Status PaymentStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := h.paymentService.UpdatePaymentStatus(paymentID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}

func (h *PaymentHandler) GetPaymentDetails(c *gin.Context) {
	paymentID := c.Param("paymentID")

	payment, err := h.paymentService.GetPaymentDetails(paymentID)
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
			_, err := h.ledgerService.RecordTransaction(c, txn)
			if err != nil {

			}
		}()
	}

	if err := h.ledgerService.UpdateTransactionStatus(c, payload.ExternalRef, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}
