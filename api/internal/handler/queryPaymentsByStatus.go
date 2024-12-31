package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. `github.com/mavrk-mose/pay/api/internal/model`
)

func (h *ApiHandler) QueryPaymentsByStatus(c *gin.Context) {
	userID := c.Param("userID")
	status := c.Query("status")

	var paymentStatus PaymentStatus
	if err := paymentStatus.UnmarshalText([]byte(status)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
		return
	}

	payments, err := h.service.QueryPaymentsByStatus(userID, paymentStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}