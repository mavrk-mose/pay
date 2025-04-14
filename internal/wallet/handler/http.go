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

// GetWalletByID returns a specific wallet by its ID
func (h *WalletHandler) GetWalletByID(c *gin.Context) {
	walletID := c.Param("walletID")

	w, err := h.service.GetWallet(c, walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": w})
}

// CreditWallet adds funds to the wallet
func (h *WalletHandler) CreditWallet(c *gin.Context) {
	var req wallet.WalletTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	err := h.service.CreditWallet(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to credit wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Credit successful"})
}

// DebitWallet deducts funds from the wallet
func (h *WalletHandler) DebitWallet(c *gin.Context) {
	var req wallet.WalletTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	err := h.service.DebitWallet(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to debit wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "Debit successful"})
}

// DeleteWallet deletes a wallet
func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	walletID := c.Param("walletID")

	err := h.service.DeleteWallet(c, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wallet deleted successfully"})
}
