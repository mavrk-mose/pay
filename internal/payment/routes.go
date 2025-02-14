package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/pkg/middleware"
	"github.com/mavrk-mose/pay/internal/payment/ports"
)

type ApiHandler struct {
	service ports.ApiService
}

func NewApiHandler(s ports.ApiService, m *middleware.ApiMiddleware, r *gin.Engine) *ApiHandler {
	handler := &ApiHandler{
		service: s,
	}

	r.GET("/check", handler.Check)

	//TODO: need to verify signature with RSA in middleware
	payment := r.Group("/api/v1", m.Authorization()) 
	{
		payment.POST("/event", handler.ReceivePaymentEvent)                                          // Receive payment event
		payment.GET("/id/:paymentID", handler.GetPaymentDetails)                                     // Get payment details
		payment.GET("/user/:userID/outgoing", handler.QueryOutgoingPayments)                         // Query outgoing payments
		payment.GET("/user/:userID/incoming", handler.QueryIncomingPayments)                         // Query incoming payments
		payment.GET("/user/:userID/date-range", handler.QueryPaymentsByDateRange)                    // Query payments by date range
		payment.GET("/user/:userID/status", handler.QueryPaymentsByStatus)                           // Query payments by status
		payment.PATCH("/id/:paymentID/status", handler.UpdatePaymentStatus)                          // Update payment status
		payment.GET("/id/:paymentID/status", handler.GetPaymentStatus)                               // Get payment status
		payment.POST("/id/:paymentID/authorize", handler.AuthorizePayment)          // Authorize payment
		payment.POST("/id/:paymentID/process", handler.ProcessPayment)              // Process authorized payment
	}

	return handler
}

