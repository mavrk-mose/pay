package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/adyen/adyen-go-api-library/v7/src/checkout"
	"github.com/adyen/adyen-go-api-library/v7/src/common"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type AdyenProvider struct {
	Client *checkout.APIClient
	MerchantAccount string
}

func NewAdyenProvider(apiKey, merchantAccount string, isLive bool) (*AdyenProvider, error) {
	env := common.TestEnv
	if isLive {
		env = common.LiveEnv
	}

	client := checkout.NewAPIClient(&checkout.Configuration{
		APIKey: apiKey,
		Env:    env,
	})

	return &AdyenProvider{Client: client, MerchantAccount: merchantAccount}, nil
}

// ExecutePayment processes a payment request using Adyen
func (a *AdyenProvider) ExecutePayment(order PaymentOrder) (*checkout.PaymentResponse, error) {
	request := checkout.PaymentRequest{
		Amount: checkout.Amount{
			Currency: order.Currency,
			Value:    int64(order.Amount * 100), // Convert to minor units
		},
		MerchantAccount: a.MerchantAccount,
		Reference:       order.OrderID,
		PaymentMethod: map[string]interface{}{
			"type":        "scheme",    // Card Payments
			"number":      order.Card.Number,
			"expiryMonth": order.Card.ExpiryMonth,
			"expiryYear":  order.Card.ExpiryYear,
			"cvc":         order.Card.CVC,
			"holderName":  order.Card.HolderName,
		},
		ReturnURL: "https://yourdomain.com/payment-result",
	}

	response, httpResp, err := a.Client.PaymentsApi.Payments(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}
	defer httpResp.Body.Close()

	return &response, nil
}

// CapturePayment captures an authorized payment
func (a *AdyenProvider) CapturePayment(paymentID string, amount float64, currency string) (*checkout.ModificationResponse, error) {
	request := checkout.CaptureRequest{
		Amount: checkout.Amount{
			Currency: currency,
			Value:    int64(amount * 100),
		},
		MerchantAccount: a.MerchantAccount,
	}

	response, httpResp, err := a.Client.PaymentsApi.Captures(context.Background(), paymentID, request)
	if err != nil {
		return nil, fmt.Errorf("capture failed: %w", err)
	}
	defer httpResp.Body.Close()

	return &response, nil
}

// RefundPayment issues a refund for a payment
func (a *AdyenProvider) RefundPayment(paymentID string, amount float64, currency string) (*checkout.ModificationResponse, error) {
	request := checkout.RefundRequest{
		Amount: checkout.Amount{
			Currency: currency,
			Value:    int64(amount * 100),
		},
		MerchantAccount: a.MerchantAccount,
	}

	response, httpResp, err := a.Client.PaymentsApi.Refunds(context.Background(), paymentID, request)
	if err != nil {
		return nil, fmt.Errorf("refund failed: %w", err)
	}
	defer httpResp.Body.Close()

	return &response, nil
}

// CancelPayment cancels an authorized payment
func (a *AdyenProvider) CancelPayment(paymentID string) (*checkout.ModificationResponse, error) {
	request := checkout.CancelRequest{
		MerchantAccount: a.MerchantAccount,
	}

	response, httpResp, err := a.Client.PaymentsApi.Cancels(context.Background(), paymentID, request)
	if err != nil {
		return nil, fmt.Errorf("cancel failed: %w", err)
	}
	defer httpResp.Body.Close()

	return &response, nil
}
