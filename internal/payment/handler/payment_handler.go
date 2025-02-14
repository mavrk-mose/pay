package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"payment-system/internal/payment/models"
	"payment-system/internal/payment/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService}
}

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