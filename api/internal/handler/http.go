package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/api/internal/middleware"
	"github.com/mavrk-mose/pay/api/internal/ports"
)

type ApiHandler struct {
	service ports.ApiService
}

func NewApiHandler(s ports.ApiService, m *middleware.ApiMiddleware, r *gin.Engine) *ApiHandler {
	handler := &ApiHandler{
		service: s,
	}

	r.GET("/check", handler.Check)

	payment := r.Group("/payment") 
	{
		payment.POST("/event", handler.ReceivePaymentEvent)                                          // Receive payment event
		payment.GET("/id/:paymentID", handler.GetPaymentDetails)                                     // Get payment details
		payment.GET("/user/:userID/outgoing", handler.QueryOutgoingPayments)                         // Query outgoing payments
		payment.GET("/user/:userID/incoming", handler.QueryIncomingPayments)                         // Query incoming payments
		payment.GET("/user/:userID/date-range", handler.QueryPaymentsByDateRange)                    // Query payments by date range
		payment.GET("/user/:userID/status", handler.QueryPaymentsByStatus)                           // Query payments by status
		payment.PATCH("/id/:paymentID/status", handler.UpdatePaymentStatus)                          // Update payment status
		payment.GET("/id/:paymentID/status", handler.GetPaymentStatus)                               // Get payment status
		payment.POST("/id/:paymentID/authorize", m.CheckAuthor(), handler.AuthorizePayment)          // Authorize payment
		payment.POST("/id/:paymentID/process", m.CheckAuthor(), handler.ProcessPayment)              // Process authorized payment
	}

	return handler
}

