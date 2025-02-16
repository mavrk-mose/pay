package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/executor"
	fraud "github.com/mavrk-mose/pay/internal/fraud/models"
	ledger "github.com/mavrk-mose/pay/internal/ledger/service"
	"github.com/mavrk-mose/pay/internal/payment/models"
	wallet "github.com/mavrk-mose/pay/internal/wallet/service"
	"net/http"
	"time"
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

func (h *PaymentService) ProcessPayment(c *gin.Context, req models.PaymentIntent) error {
	// 1. Check balance
	balance, err := h.walletService.GetBalance(c, req.ID)
	if err != nil || balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return err
	}

	// 2. Execute payment
	err, _ = h.executor.ExecutePayment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return err
	}

	// 3. Record transaction in ledger
	txn := fraud.Transaction{
		ExternalRef:    req.ReceiptNumber,
		Type:           fraud.TransactionType("deposit"),
		Status:         fraud.TransactionStatus("pending"),
		Details:        req.Description,
		Currency:       req.Currency,
		DebitWalletID:  1121213123, //TODO: rewrite this
		DebitAmount:    float64(req.Amount),
		EntryType:      fraud.EntryType("debit"),
		CreditWalletID: 213412434354, // TODO: rewrite this
		CreditAmount:   float64(req.Amount),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	h.ledgerService.RecordTransaction(txn)

	// 4. Deduct from wallet
	h.walletService.UpdateBalance(req.UserID, -req.Amount)

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}
