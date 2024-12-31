package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ApiHandler) AuthorizePayment(c *gin.Context) {
	paymentID := c.Param("paymentID")

	authorized, err := h.service.AuthorizePayment(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorized": authorized})
}