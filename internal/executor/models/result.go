package models

type PaymentStatus string

const (
	StatusSuccess PaymentStatus = "success"
	StatusFailed  PaymentStatus = "failed"
	StatusPending PaymentStatus = "pending"
)

// PaymentResult encapsulates a generic response from any payment provider.
type PaymentResult struct {
	Status         PaymentStatus 
	Message        string                 
	TransactionID  string                 
	Provider       string                 
	ProviderRefID  string                 
	RawResponse    any                    
	ErrorCode      string                 
	ErrorDetail    string                 
}

func (r PaymentResult) Error() string {
	if r.Status != StatusSuccess {
		return r.Message + " (" + r.ErrorCode + ")"
	}
	return ""
}
