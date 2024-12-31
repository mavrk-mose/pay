package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/api/internal/model"
)

func (h *ApiHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")
	var request struct {
		Status PaymentStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := h.service.UpdatePaymentStatus(paymentID, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment status updated successfully"})
}