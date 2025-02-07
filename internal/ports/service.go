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