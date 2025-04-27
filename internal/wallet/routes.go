package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	walletHandlers "github.com/mavrk-mose/pay/internal/wallet/handler"
	"github.com/mavrk-mose/pay/internal/wallet/repository"
	"github.com/mavrk-mose/pay/internal/wallet/service"
	"github.com/mavrk-mose/pay/pkg/middleware"
)

func NewApiHandler(r *gin.Engine, db *sqlx.DB) {
	walletService := service.NewWalletService(repository.NewWalletRepo(db))
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
