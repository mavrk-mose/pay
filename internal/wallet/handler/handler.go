package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/user/repository"
)

type WalletHandler struct {
	repo repository.UserRepository
}

func NewWalletHandler(repo repository.UserRepository) *WalletHandler {
	return &WalletHandler{repo: repo}
}

// CreateWallet allows a user to create a new wallet
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID := c.Param("userID")
	var req struct {
		Currency string `json:"currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	//TODO: move the create wallet from user repo to wallet repo
	wallet, err := h.repo.CreateWallet(userID, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}

// GetUserWallets retrieves all wallets for a user
func (h *WalletHandler) GetUserWallets(c *gin.Context) {
	userID := c.Param("userID")

	wallets, err := h.repo.GetUserWallets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wallets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}
