package handler

import (
	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"net/http"
)

func (t *PaymentHandler) ProcessPayment(ctx *gin.Context) {
	var paymentIntent PaymentIntent

	if err := ctx.ShouldBindJSON(&paymentIntent); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := t.service.ProcessPayment(ctx, paymentIntent); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully"})
}
