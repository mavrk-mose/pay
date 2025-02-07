package services

type AdyenGateway struct{}

func (a *AdyenGateway) ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error) {
	// Simulate Adyen API call.
	return PaymentExecutionResult{
		Success:       true,
		Message:       "Payment processed by Adyen",
		TransactionID: "adyen_txn_abcde",
	}, nil
}