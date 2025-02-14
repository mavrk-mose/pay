package service

import (
	"github.com/gin-gonic/gin"
	executor "github.com/mavrk-mose/pay/internal/executor"
	ledger "github.com/mavrk-mose/pay/internal/ledger/service"
	wallet "github.com/mavrk-mose/pay/internal/wallet/service"
	"net/http"
)

type PaymentService struct {
	walletService *wallet.WalletService
	ledgerService *ledger.LedgerService
	executor      *executor.PaymentExecutor
}

func NewPaymentService(wallet *wallet.WalletService, ledger *ledger.LedgerService, executor *executor.PaymentExecutor) *PaymentService {
	return &PaymentService{wallet, ledger, executor}
}

type ExternalPaymentResponse struct {
	Status      string `json:"status"`
	ExternalRef string `json:"external_ref"`
}

func (h *PaymentService) ProcessPayment(c *gin.Context) {
	var req struct {
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		Gateway string  `json:"gateway"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Check balance
	balance, err := h.walletService.GetBalance(req.UserID)
	if err != nil || balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// 2. Execute payment
	err = h.executor.ExecutePayment(req.Amount, req.UserID, req.Gateway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return
	}

	// 3. Record transaction in ledger
	h.ledgerService.RecordTransaction(req.UserID, "DEBIT", req.Amount)

	// 4. Deduct from wallet
	h.walletService.UpdateBalance(req.UserID, -req.Amount)

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}
