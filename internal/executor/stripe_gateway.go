package executor

import (
	. "github.com/mavrk-mose/pay/internal/payment/models"
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
