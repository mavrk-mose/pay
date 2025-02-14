package ports

import (
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"time"
)

type ApiService interface {
	ReceivePaymentEvent(event PaymentEvent) error
	QueryOutgoingPayments(userID string) ([]PaymentOrder, error)
	QueryIncomingPayments(userID string) ([]PaymentOrder, error)
	AuthorizePayment(paymentID string) (bool, error) // Checks if a payment can be authorized
	ProcessPayment(paymentID string) error           // Processes an authorized payment

	// user requests
	GetPaymentDetails(paymentID string) (PaymentOrder, error)                                     // Retrieves detailed information about a payment
	QueryPaymentsByDateRange(userID string, startDate, endDate time.Time) ([]PaymentOrder, error) // Payments by date
	QueryPaymentsByStatus(userID string, status PaymentStatus) ([]PaymentOrder, error)            // Payments by status

	// status checks
	UpdatePaymentStatus(paymentID string, status PaymentStatus) error // Updates the status of a payment
	GetPaymentStatus(paymentID string) (PaymentStatus, error)         // Retrieves the current status of a payment
}

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
