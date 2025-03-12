package ports

import (
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"time"
)

type ApiStore interface {
	SavePayment(payment PaymentOrder) error
	UpdatePaymentStatus(paymentID string, status PaymentStatus) error
	GetPaymentByID(paymentID string) (PaymentOrder, error)
	GetOutgoingPayments(userID string) ([]PaymentOrder, error)
	GetIncomingPayments(userID string) ([]PaymentOrder, error)
	GetPaymentsByDateRange(userID string, startDate, endDate time.Time) ([]PaymentOrder, error)
	GetPaymentsByStatus(userID string, status PaymentStatus) ([]PaymentOrder, error)
	DeletePayment(paymentID string) error
}
