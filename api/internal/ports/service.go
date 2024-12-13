package ports

import . "github.com/mavrk-mose/pay/api/internal/model"

type ApiService interface {
	ReceivePaymentEvent(event PaymentEvent) error
	QueryOutgoingPayments(userID string) ([]PaymentOrder, error)
	QueryIncomingPayments(userID string) ([]PaymentOrder, error)
}
