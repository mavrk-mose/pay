package executor

type PaymentExecutionResult struct {
	Success       bool
	Message       string
	TransactionID string
}
