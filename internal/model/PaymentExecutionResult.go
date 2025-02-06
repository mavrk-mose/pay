package model

type PaymentExecutionResult struct {
	OrderID     string
	Status      string // e.g., "Success", "Failed" -> can be enum of status codes
	PSPResponse string // Response from payment service provider
}
