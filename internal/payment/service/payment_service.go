package service

import (
	"encoding/json"
	"errors"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/mavrk-mose/pay/internal/payment/repository"
	"github.com/mavrk-mose/pay/internal/wallet"
	"github.com/mavrk-mose/pay/internal/ledger"
	"github.com/mavrk-mose/pay/internal/executor"
	"github.com/mavrk-mose/pay/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type PaymentActions interface {
	RecordPaymentEvent(event PaymentEvent) error
	ForwardPaymentOrder(order PaymentOrder) error
	NotifyLedger(order PaymentOrder) error
	NotifyWallet(order PaymentOrder) error

	SchedulePayment(order PaymentOrder, scheduleTime time.Time) error // Schedules a payment for future processing
	GetScheduledPayments(userID string) ([]PaymentOrder, error)       // Retrieves scheduled payments
	CancelScheduledPayment(paymentID string) error                    // Cancels a scheduled payment

	LogFailedPayment(event PaymentEvent, reason string) error // Logs a failed payment event
	RetryFailedPayment(paymentID string) error                // Retries processing a failed payment

	GetTransactionVolumeByPeriod(period string) (int, error) // Retrieves total transaction volume by period

	FlagSuspiciousPayment(paymentID string) error               // Flags a payment as suspicious
	UnflagPayment(paymentID string) error                       // Removes suspicious flag from a payment
	QueryFlaggedPayments(userID string) ([]PaymentOrder, error) // Lists flagged payments

	RequestRefund(paymentID string) error     // Initiates a refund for a payment
	CancelPayment(paymentID string) error     // Cancels a pending payment
	QueryRefundStatus(paymentID string) error // Queries the status of a refund
}

type PaymentService struct {
	walletService  *wallet.WalletService
	ledgerService  *ledger.LedgerService
	executor       *executor.PaymentExecutor
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


// type PaymentService struct {
// 	TransactionRepo    repository.TransactionRepo
// 	HttpClient         *http.Client
// 	ExternalPaymentURL string
// 	GenericHttpClient  *utils.GenericHttpClient
// }

// func NewPaymentService(transactionRepo repository.TransactionRepo, externalPaymentURL string, logger *zap.Logger) *PaymentService {
// 	return &PaymentService{
// 		TransactionRepo:    transactionRepo,
// 		ExternalPaymentURL: externalPaymentURL,
// 		GenericHttpClient:  utils.NewGenericHttpClient(logger),
// 	}
// }

// func (s *PaymentService) InitiatePayment(paymentRequest PaymentIntent) (string, error) {
// 	extResp, err := s.GenericHttpClient.Post(s.ExternalPaymentURL, paymentRequest, map[string]string{})
// 	if err != nil {
// 		return "", err
// 	}

// 	var response ExternalPaymentResponse
// 	if err := json.Unmarshal(*extResp, &response); err != nil {
// 		return "", err
// 	}

// 	if response.Status != "ok" && response.Status != "confirmed" {
// 		return "", errors.New("external payment initiation failed")
// 	}

// 	return response.ExternalRef, nil
// }