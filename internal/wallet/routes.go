package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/config"
	walletHandlers "github.com/mavrk-mose/pay/internal/wallet/handler"
	"github.com/mavrk-mose/pay/internal/wallet/service"
	"github.com/mavrk-mose/pay/pkg/middleware"
)

func NewApiHandler(r *gin.Engine, cfg *config.Config) {
	walletService := service.NewWalletService()
	walletHandler := walletHandlers.NewWalletHandler(walletService)

	api := r.Group("/api/v1", middleware.AuthMiddleware())

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
}
