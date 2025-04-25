package service

import (
	"errors"
	. "github.com/mavrk-mose/pay/internal/executor/models"
	models "github.com/mavrk-mose/pay/internal/payment/models"
)

type PaymentExecutorService interface {
	ExecutePayment(order models.PaymentIntent) (any, error)
	RecordPaymentOrder(order models.PaymentIntent) error
}

type PaymentGateway interface {
	ExecutePayment(order models.PaymentOrder) (any, error)
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

func (pe *PaymentExecutor) ExecutePayment(order models.PaymentIntent) (any, error) {
	gateway, exists := pe.gateways[order.PaymentMethod]
	if !exists {
		return PaymentResult{}, errors.New("unsupported payment gateway: " + order.PaymentMethod)
	}
	paymentOrder := models.PaymentOrder{
		Amount:        order.Amount,
		Currency:      order.Currency,
		Description:   order.Description,
		PayerID:       order.Customer.String(),
		PayeeID:       order.ReceiptNumber,
		PaymentMethod: order.PaymentMethod,
	}
	return gateway.ExecutePayment(paymentOrder)
}

func (pe *PaymentExecutor) RecordPaymentOrder(order models.PaymentIntent) error {
	// Implementation for recording payment order
	return nil
}
