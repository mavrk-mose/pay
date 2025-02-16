package handler

import (
	"github.com/mavrk-mose/pay/internal/payment/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")
	var request struct {
		Status models.PaymentStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := t.service.UpdatePaymentStatus(paymentID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}
