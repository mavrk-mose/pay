package ports

import (
	"time"

	. "github.com/mavrk-mose/pay/api/internal/model"
)

type ApiService interface {
	ReceivePaymentEvent(event PaymentEvent) error
	QueryOutgoingPayments(userID string) ([]PaymentOrder, error)
	QueryIncomingPayments(userID string) ([]PaymentOrder, error)
	AuthorizePayment(paymentID string) (bool, error) // Checks if a payment can be authorized
	ProcessPayment(paymentID string) error          // Processes an authorized payment

	// user requests
	GetPaymentDetails(paymentID string) (PaymentOrder, error) // Retrieves detailed information about a payment
	QueryPaymentsByDateRange(userID string, startDate, endDate time.Time) ([]PaymentOrder, error) // Payments by date
	QueryPaymentsByStatus(userID string, status PaymentStatus) ([]PaymentOrder, error) // Payments by status

	// status checks
	UpdatePaymentStatus(paymentID string, status PaymentStatus) error // Updates the status of a payment
	GetPaymentStatus(paymentID string) (PaymentStatus, error)         // Retrieves the current status of a payment
}

type PaymentExecutorService interface {
	ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error)
	RecordPaymentOrder(order PaymentOrder) error
}

type LedgerService interface {
	RecordTransaction(transaction LedgerTransaction) error
	GetTransactionByID(transactionID string) (LedgerTransaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}

type ReconciliationService interface {
	ReconcileLedgerWithSettlement(ledger LedgerService, settlementFile SettlementFile) (ReconciliationResult, error)
}

type PaymentService interface {
	RecordPaymentEvent(event PaymentEvent) error
	ForwardPaymentOrder(order PaymentOrder) error
	NotifyLedger(order PaymentOrder) error
	NotifyWallet(order PaymentOrder) error

	//scheduled payments
	SchedulePayment(order PaymentOrder, scheduleTime time.Time) error // Schedules a payment for future processing
	GetScheduledPayments(userID string) ([]PaymentOrder, error)       // Retrieves scheduled payments
	CancelScheduledPayment(paymentID string) error                    // Cancels a scheduled payment

	//error handling & logging
	LogFailedPayment(event PaymentEvent, reason string) error // Logs a failed payment event
	RetryFailedPayment(paymentID string) error                // Retries processing a failed payment

	//analytics
	// GetPaymentSummary(userID string) (PaymentSummary, error) // Provides a summary of user's payments
	GetTransactionVolumeByPeriod(period string) (int, error) // Retrieves total transaction volume by period

	// fraud detetion
	FlagSuspiciousPayment(paymentID string) error      // Flags a payment as suspicious
	UnflagPayment(paymentID string) error             // Removes suspicious flag from a payment
	QueryFlaggedPayments(userID string) ([]PaymentOrder, error) // Lists flagged payments

	//cancel & refunds
	RequestRefund(paymentID string) error      // Initiates a refund for a payment
	CancelPayment(paymentID string) error      // Cancels a pending payment
	QueryRefundStatus(paymentID string) error // Queries the status of a refund
}

type WalletService interface {
	UpdateBalance(userID string, amount float64) error
	GetBalance(userID string) (float64, error)
}