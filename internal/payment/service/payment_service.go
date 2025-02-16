package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/executor"
	. "github.com/mavrk-mose/pay/internal/fraud/models"
	ledger "github.com/mavrk-mose/pay/internal/ledger/service"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	wallet "github.com/mavrk-mose/pay/internal/wallet/service"
)

// PaymentService handles payment processing
type PaymentService struct {
	walletService wallet.WalletService
	ledgerService ledger.LedgerService
	executor      executor.PaymentExecutor
}

func NewPaymentService(wallet wallet.WalletService, ledger ledger.LedgerService, executor executor.PaymentExecutor) *PaymentService {
	return &PaymentService{
		walletService: wallet,
		ledgerService: ledger,
		executor:      executor,
	}
}

// ExternalPaymentResponse represents a response from an external payment provider
type ExternalPaymentResponse struct {
	Status      string `json:"status"`
	ExternalRef string `json:"external_ref"`
}

// ProcessPayment handles the entire payment flow
func (h *PaymentService) ProcessPayment(ctx *gin.Context, req PaymentIntent) error {
	// 1️⃣ Check wallet balance
	balance, err := h.walletService.GetBalance(ctx, req.Customer)
	if err != nil || balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return err
	}

	// 2️⃣ Execute payment via external processor
	err, _ = h.executor.ExecutePayment(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		return err
	}

	// 3️⃣ Record transaction in ledger
	txn := Transaction{
		ExternalRef:    req.ReceiptNumber,
		Type:           TransactionType("payment"),
		Status:         TransactionStatus("pending"),
		Details:        req.Description,
		Currency:       req.Currency,
		DebitWalletID:  324532453245, // Use actual wallet ID
		DebitAmount:    req.Amount,
		EntryType:      EntryType("debit"),
		CreditWalletID: 235432453455, // Use actual recipient wallet ID
		CreditAmount:   req.Amount,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := h.ledgerService.RecordTransaction(ctx, txn); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction in ledger"})
		return err
	}

	// 4️⃣ Deduct funds from sender's wallet
	if err := h.walletService.UpdateBalance(ctx, req.Customer, -req.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return err
	}

	// ✅ Payment Successful
	ctx.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
	return nil
}

func (h *PaymentService) GetPaymentDetails(id string) (PaymentIntent, error) {
	return PaymentIntent{}, nil
}

func (h *PaymentService) GetPaymentStatus(id string) (PaymentStatus, error) {
	return "", nil
}

func (h *PaymentService) QueryIncomingPayments(id string) ([]PaymentIntent, error) {
	return []PaymentIntent{}, nil
}

func (h *PaymentService) QueryOutgoingPayments(id string) ([]PaymentIntent, error) {
	return []PaymentIntent{}, nil
}

func (h *PaymentService) QueryPaymentsByDateRange(id string, date time.Time, date2 time.Time) ([]PaymentIntent, error) {
	return []PaymentIntent{}, nil
}

func (h *PaymentService) QueryPaymentsByStatus(id string, status PaymentStatus) ([]PaymentIntent, error) {
	return []PaymentIntent{}, nil
}

func (h *PaymentService) UpdatePaymentStatus(id string, status PaymentStatus) error {
	return nil
}
