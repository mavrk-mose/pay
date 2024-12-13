package ports

import . "github.com/mavrk-mose/pay/executor/internal/model"

type PaymentExecutorService interface {
	ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error)
	RecordPaymentOrder(order PaymentOrder) error
}
