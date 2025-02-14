package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) QueryPaymentsByDateRange(c *gin.Context) {
	userID := c.Param("userID")
	startDate, err1 := time.Parse(time.RFC3339, c.Query("startDate"))
	endDate, err2 := time.Parse(time.RFC3339, c.Query("endDate"))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
		return
	}

	payments, err := h.service.QueryPaymentsByDateRange(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}