package service

import (
	"errors"
	. "github.com/mavrk-mose/pay/internal/executor/models"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type PaymentExecutorService interface {
	ExecutePayment(order PaymentIntent) (PaymentResult, error)
	RecordPaymentOrder(order PaymentIntent) error
}

type PaymentGateway interface {
	ExecutePayment(order PaymentOrder) (PaymentResult, error)
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
		"stripe": &StripeProvider{},
		"paypal": &PayPalProvider{},
		"adyen":  &AdyenProvider{},
	}
	return NewPaymentExecutor(gateways)
}

func (pe *PaymentExecutor) ExecutePayment(order PaymentIntent) (PaymentResult, error) {
	gateway, exists := pe.gateways[order.PaymentMethod]
	if !exists {
		return PaymentResult{}, errors.New("unsupported payment gateway: " + order.PaymentMethod)
	}
	paymentOrder := PaymentOrder{
		Amount:        order.Amount,
		Currency:      order.Currency,
		Description:   order.Description,
		PayerID:       order.Customer.String(),
		PayeeID:       order.ReceiptNumber,
		PaymentMethod: order.PaymentMethod,
	}
	return gateway.ExecutePayment(paymentOrder)
}

func (pe *PaymentExecutor) RecordPaymentOrder(order PaymentIntent) error {
	// Implementation for recording payment order
	return nil
}
