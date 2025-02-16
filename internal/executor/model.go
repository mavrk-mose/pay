package executor

type PaymentExecutionResult struct {
	Success       bool
	Message       string
	TransactionID string
}

func (p PaymentExecutionResult) Error() string {
	//TODO implement me
	panic("implement me")
}
