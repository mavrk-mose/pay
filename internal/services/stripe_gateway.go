package services

import (
	. "github.com/mavrk-mose/pay/internal/model"
)

type StripeGateway struct{}

func (s *StripeGateway) ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error) {
	// Simulate Stripe API call.
	return PaymentExecutionResult{
		Success:       true,
		Message:       "Payment processed by Stripe",
		TransactionID: "stripe_txn_12345",
	}, nil
}
