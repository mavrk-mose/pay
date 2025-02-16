package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/payment/models"
)

// QueryPayments allows filtering payments based on type, date range, and status
func (t *PaymentHandler) QueryPayments(c *gin.Context) {
	userID := c.Param("userID")
	paymentType := c.Query("type") // "incoming" or "outgoing"
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
	var paymentStatus models.PaymentStatus
	if status != "" {
		if err := paymentStatus.UnmarshalText([]byte(status)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
			return
		}
	}

	// Query payments based on filters
	var payments []models.PaymentIntent
	switch paymentType {
	case "incoming":
		payments, err = t.service.QueryIncomingPayments(userID)
	case "outgoing":
		payments, err = t.service.QueryOutgoingPayments(userID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter. Use 'incoming' or 'outgoing'"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by date range if applicable
	if !startDate.IsZero() && !endDate.IsZero() {
		payments, err = t.service.QueryPaymentsByDateRange(userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Filter by status if applicable
	if status != "" {
		payments, err = t.service.QueryPaymentsByStatus(userID, paymentStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return filtered payments
	c.JSON(http.StatusOK, payments)
}
