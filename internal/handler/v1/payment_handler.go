package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/api/internal/model"
)

func (h *ApiHandler) ReceivePaymentEvent(c *gin.Context) {
	var event PaymentIntent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := h.service.ReceivePaymentEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment event received successfully"})
}