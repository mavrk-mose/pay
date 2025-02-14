package executor

import (
	"errors"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

// Payment execution (Stripe, PayPal, Bank API)

type PaymentExecutorService interface {
	ExecutePayment(order PaymentIntent) (PaymentExecutionResult, error)
	RecordPaymentOrder(order PaymentIntent) error
}

type PaymentGateway interface {
	ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error)
}

type PaymentExecutor struct {
	gateways map[string]PaymentGateway
}

func NewPaymentExecutor(gateways map[string]PaymentGateway) PaymentExecutorService {
	return &PaymentExecutor{
		gateways: gateways,
	}
}

func NewDefaultPaymentExecutor() PaymentExecutorService {
	gateways := map[string]PaymentGateway{
		"stripe": &StripeGateway{},
		"paypal": &PaypalGateway{},
		"adyen":  &AdyenGateway{},
	}
	return NewPaymentExecutor(gateways)
}

func (pe *PaymentExecutor) ExecutePayment(order PaymentIntent) (PaymentExecutionResult, error) {
	gateway, exists := pe.gateways[order.PaymentMethod]
	if !exists {
		return PaymentExecutionResult{}, errors.New("unsupported payment gateway: " + order.PaymentMethod)
	}
	paymentOrder := PaymentOrder{
		Amount:        float64(order.Amount),
		Currency:      order.Currency,
		Description:   order.Description,
		PayerID:       order.Customer,
		PayeeID:       order.ReceiptNumber,
		PaymentMethod: order.PaymentMethod,
	}
	return gateway.ExecutePayment(paymentOrder)
}

func (pe *PaymentExecutor) RecordPaymentOrder(order PaymentIntent) error {
	// Implementation for recording payment order
	return nil
}

// type PaymentService struct {
// 	TransactionRepo    repository.TransactionRepo
// 	HttpClient         *http.Client
// 	ExternalPaymentURL string
// 	GenericHttpClient  *utils.GenericHttpClient
// }

// func NewPaymentService(transactionRepo repository.TransactionRepo, externalPaymentURL string, logger *zap.Logger) *PaymentService {
// 	return &PaymentService{
// 		TransactionRepo:    transactionRepo,
// 		ExternalPaymentURL: externalPaymentURL,
// 		GenericHttpClient:  utils.NewGenericHttpClient(logger),
// 	}
// }

// func (s *PaymentService) InitiatePayment(paymentRequest PaymentIntent) (string, error) {
// 	extResp, err := s.GenericHttpClient.Post(s.ExternalPaymentURL, paymentRequest, map[string]string{})
// 	if err != nil {
// 		return "", err
// 	}

// 	var response ExternalPaymentResponse
// 	if err := json.Unmarshal(*extResp, &response); err != nil {
// 		return "", err
// 	}

// 	if response.Status != "ok" && response.Status != "confirmed" {
// 		return "", errors.New("external payment initiation failed")
// 	}

// 	return response.ExternalRef, nil
// }
