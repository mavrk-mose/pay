package handler

import (
	wallet "github.com/mavrk-mose/pay/internal/wallet/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/wallet/service"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

// CreateWallet allows a user to create a new wallet
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req wallet.CreateWalletRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	dbWallet, err := h.service.CreateWallet(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": dbWallet})
}

// GetUserWallets retrieves all wallets for a user
func (h *WalletHandler) GetUserWallets(c *gin.Context) {
	userID := c.Param("userID")

	wallets, err := h.service.GetUserWallets(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}
