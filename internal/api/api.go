package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/api/middleware"
	v1 "github.com/mavrk-mose/pay/internal/api/v1"
)


func NewApiHandler(r *gin.Engine, db *sqlx.DB, cfg *config.Config) {
	
	walletHandler := v1.NewWalletHandler(db)
	userHandler := v1.NewUserHandler(db)
	paymentHandler := v1.NewPaymentHandler(db)
	notificationHandler := v1.NewNotificationHandler(cfg, db)
	executorHandler := v1.NewWebhookHandler()

	api := r.Group("/api/v1", middleware.AuthMiddleware())

	v1.InitAuth(cfg)

	r.GET("/check", paymentHandler.Check)

	// Authentication routes
	auth := api.Group("/auth") 
	{
		auth.GET("/:provider", v1.BeginAuthHandler)
		auth.GET("/:provider/callback", userHandler.AuthCallbackHandler)
		auth.GET("/logout/:provider", userHandler.LogoutHandler)
	}

	// Admin routes
	admin := api.Group("/admin/users", middleware.AdminMiddleware(cfg))
	{
		admin.GET("/", userHandler.ListUsers)
		admin.POST("/:userID/assign-role", userHandler.AssignRole)
		admin.POST("/:userID/revoke-role", userHandler.RevokeRole)
		admin.POST("/:userID/ban", userHandler.BanUser)
		admin.POST("/:userID/unban", userHandler.UnbanUser)
	}

	// Wallet Routes
	wallet := api.Group("/wallet")
	{
		wallet.POST("/", walletHandler.CreateWallet)
		wallet.GET("/user/:userID", walletHandler.GetUserWallets)
		wallet.GET("/:walletID", walletHandler.GetWalletByID)
		wallet.POST("/credit", walletHandler.CreditWallet)
		wallet.POST("/debit", walletHandler.DebitWallet)
		wallet.DELETE("/:walletID", walletHandler.DeleteWallet)
	}

	// Payment routes
	payment := api.Group("/payment")
	{
		payment.POST("/event", paymentHandler.ProcessPayment)                      
		payment.GET("/id/:paymentID", paymentHandler.GetPaymentDetails)            
		payment.GET("/user/:userID/date-range", paymentHandler.QueryPayments)     
		payment.GET("/user/:userID/status", paymentHandler.QueryPayments)          
		payment.PATCH("/id/:paymentID/status", paymentHandler.UpdatePaymentStatus) 
		payment.GET("/id/:paymentID/status", paymentHandler.GetPaymentStatus)     
		payment.POST("/id/:paymentID", paymentHandler.GetPaymentDetails)           
		payment.POST("/id/:paymentID/process", paymentHandler.ProcessPayment) 
	}

	notification := api.Group("/notifications")
	{
		notification.GET("/", notificationHandler.GetNotifications) // should be realtime -> ideally websockets
		notification.POST("/:id/read", notificationHandler.MarkAsRead)
	}

	executor := api.Group("/webhook") //validate with signature
	{
		executor.POST("/stripe", executorHandler.StripeWebhookHandler)
		executor.POST("/paypal", executorHandler.PaypalWebhookHandler)
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