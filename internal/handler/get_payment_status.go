package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("paymentID")

	status, err := h.service.GetPaymentStatus(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}