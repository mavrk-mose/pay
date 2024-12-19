package ports

import (
	"time"

	. "github.com/mavrk-mose/pay/service/internal/model"
)

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
