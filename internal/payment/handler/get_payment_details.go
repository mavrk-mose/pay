package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *PaymentHandler) GetPaymentDetails(c *gin.Context) {
	paymentID := c.Param("paymentID")

	payment, err := t.service.GetPaymentDetails(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
