package executor

import (
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type AdyenGateway struct{}

func (a *AdyenGateway) ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error) {
	// Simulate Adyen API call.
	return PaymentExecutionResult{
		Success:       true,
		Message:       "Payment processed by Adyen",
		TransactionID: "adyen_txn_abcde",
	}, nil
}
