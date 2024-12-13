package ports

import (
	. "github.com/mavrk-mose/pay/service/internal/model"
)

type PaymentService interface {
	RecordPaymentEvent(event PaymentEvent) error
	ForwardPaymentOrder(order PaymentOrder) error
	NotifyLedger(order PaymentOrder) error
	NotifyWallet(order PaymentOrder) error
}
