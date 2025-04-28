package service

import (
	. "github.com/mavrk-mose/pay/internal/ledger/models"
	"github.com/mavrk-mose/pay/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	executor "github.com/mavrk-mose/pay/internal/executor/service"
	ledger "github.com/mavrk-mose/pay/internal/ledger/service"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	. "github.com/mavrk-mose/pay/internal/payment/repository"
	wallet "github.com/mavrk-mose/pay/internal/wallet/service"
)

type PaymentService struct {
	walletService     wallet.WalletService
	ledgerService     ledger.LedgerService
	executor          executor.PaymentExecutor
	productConfigRepo ProductConfigRepo
	logger            utils.Logger
}

func NewPaymentService(
	wallet wallet.WalletService,
	ledger ledger.LedgerService,
	executor executor.PaymentExecutor,
	productConfigRepo ProductConfigRepo,
) *PaymentService {
	return &PaymentService{
		walletService:     wallet,
		ledgerService:     ledger,
		executor:          executor,
		productConfigRepo: productConfigRepo,
	}
}

type ExternalPaymentResponse struct {
	Status      string `json:"status"`
	ExternalRef string `json:"external_ref"`
}

func (h *PaymentService) ProcessPayment(ctx *gin.Context, req PaymentIntent) error {
	productConfig, err := h.productConfigRepo.GetProductConfig(ctx, req.ProductName)
	if err != nil {
		h.logger.Errorf("Failed to fetch product configuration: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product configuration"})
		panic(err)
	}

	balance, err := h.walletService.GetBalance(ctx, req.Customer)
	if err != nil || balance < req.Amount {
		h.logger.Errorf("Insufficient balance: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		panic(err)
	}

	feeAmount := req.Amount * productConfig.FeePercentage / 100
	netAmount := req.Amount - feeAmount

	_, err = h.executor.ExecutePayment(req)
	if err != nil {
		h.logger.Errorf("Payment failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Payment failed"})
		panic(err)
	}

	txn := Transaction{
		ExternalRef:    req.ReceiptNumber,
		Type:           TransactionType("payment"),
		Status:         TransactionStatus("pending"),
		Details:        req.Description,
		Currency:       req.Currency,
		DebitWalletID:  324532453245, // Use actual wallet ID
		Amount:         req.Amount,   // Full amount deducted from customer
		EntryType:      EntryType("debit"),
		CreditWalletID: 235432453455, // Use actual recipient wallet ID
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	go func() {
		if err := h.ledgerService.RecordTransaction(ctx, txn); err != nil {
			h.logger.Errorf("Failed to record transaction in ledger: %v", err)
			panic(err)
		}
	}()

	go func() {
		feeTxn := Transaction{
			ExternalRef:    req.ReceiptNumber + "-fee",
			Type:           TransactionType("fee"),
			Status:         TransactionStatus("completed"),
			Details:        "Transaction fee for payment",
			Currency:       req.Currency,
			DebitWalletID:  23423342424, // Customer's wallet ID
			Amount:    feeAmount,   // Fee amount deducted from customer
			EntryType:      EntryType("debit"),
			CreditWalletID: 25234534254, // System's fee wallet ID
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := h.ledgerService.RecordTransaction(ctx, feeTxn); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record fee transaction in ledger"})
			panic(err)
		}
	}()

	if err := h.walletService.UpdateBalance(ctx, req.Customer, -req.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return err
	}

	if err := h.walletService.UpdateBalance(ctx, req.Recipient, netAmount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipient's wallet balance"})
		return err
	}

	if err := h.walletService.UpdateBalance(ctx, productConfig.FeeWalletID, feeAmount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fee wallet balance"})
		return err
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
	return nil
}

func (h *PaymentService) GetPaymentDetails(id string) (PaymentIntent, error) {
	return PaymentIntent{}, nil
}

func (h *PaymentService) GetPaymentStatus(id string) (PaymentStatus, error) {
	return "", nil
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
