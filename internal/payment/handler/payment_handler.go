package handler

import (
	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"net/http"
)

func (h *PaymentHandler) ReceivePaymentEvent(c *gin.Context) {
	var event PaymentIntent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := h.paymentService.ProcessPayment(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}
