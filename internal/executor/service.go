package executor

import (
	"errors"
	. "github.com/mavrk-mose/pay/internal/model"
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
