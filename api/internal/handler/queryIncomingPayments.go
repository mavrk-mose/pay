package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) QueryIncomingPayments(c *gin.Context) {
	userID := c.Param("userID")

	payments, err := h.service.QueryIncomingPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}