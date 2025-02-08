package services

import (
	. "github.com/mavrk-mose/pay/internal/model"
)

type PaypalGateway struct{}

func (p *PaypalGateway) ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error) {
	// Simulate PayPal API call.
	return PaymentExecutionResult{
		Success:       true,
		Message:       "Payment processed by PayPal",
		TransactionID: "paypal_txn_67890",
	}, nil
}
