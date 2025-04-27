package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	repository "github.com/mavrk-mose/pay/internal/ledger/repository"
	handlers "github.com/mavrk-mose/pay/internal/payment/handler"
)

func NewApiHandler(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	handler := handlers.NewPaymentHandler(repository.Repo{DB: db})

	//publicKey, err := middleware.LoadPublicKey(cfg.Server.PublicKeyPath)
	//if err != nil {
	//	panic(err)
	//}

	r.GET("/check", handler.Check)

	payment := r.Group("/api/v1")
	{
		payment.POST("/event", handler.ProcessPayment)                      // Receive payment event
		payment.GET("/id/:paymentID", handler.GetPaymentDetails)            // Get payment details
		payment.GET("/user/:userID/date-range", handler.QueryPayments)      // Query payments by date range
		payment.GET("/user/:userID/status", handler.QueryPayments)          // Query payments by status
		payment.PATCH("/id/:paymentID/status", handler.UpdatePaymentStatus) // Update payment status
		payment.GET("/id/:paymentID/status", handler.GetPaymentStatus)      // Get payment status
		payment.POST("/id/:paymentID", handler.GetPaymentDetails)           // Authorize payment
		payment.POST("/id/:paymentID/process", handler.ProcessPayment)      // Process authorized payment
	}

	// // Referral routes
	// router.GET("/users/:userID/referral", gHandler.GenerateReferral)
	// router.POST("/users/:userID/referral", gHandler.ApplyReferral)

	// // Cashback route (query parameter: amount)
	// router.GET("/cashback", gHandler.GetCashback)

	// // Voucher routes
	// router.POST("/users/:userID/voucher", gHandler.GenerateVoucher)
	// router.POST("/users/:userID/voucher/redeem", gHandler.RedeemVoucher)

	// // Merchant discount route
	// router.GET("/merchants/:merchantID/discount", gHandler.GetMerchantDiscount)

	// // Challenges route
	// router.GET("/users/:userID/challenges", gHandler.GetChallenges)

	// // Loyalty points route
	// router.GET("/users/:userID/loyalty", gHandler.GetLoyaltyPoints)
}
