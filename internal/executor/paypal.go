package executor

import (
	"context"
	"fmt"
	"log"

	"github.com/netlify/PayPal-Go-SDK/paypal"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

// PayPalProvider struct to manage PayPal payments
type PayPalProvider struct {
	Client *paypal.Client
}

// NewPayPalProvider initializes a new PayPal client
func NewPayPalProvider(clientID, secret string, isSandbox bool) (*PayPalProvider, error) {
	var apiBase string
	if isSandbox {
		apiBase = paypal.APIBaseSandBox
	} else {
		apiBase = paypal.APIBaseLive
	}

	client, err := paypal.NewClient(clientID, secret, apiBase)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal client: %w", err)
	}

	// Automatically retrieve an access token
	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve PayPal access token: %w", err)
	}

	return &PayPalProvider{Client: client}, nil
}

// CreateOrder creates a PayPal order
func (p *PayPalProvider) CreateOrder(amount string, currency string, returnURL string, cancelURL string) (*paypal.Order, error) {
	orderIntent := paypal.IntentCapture
	orderRequest := paypal.OrderRequest{
		Intent: orderIntent,
		PurchaseUnits: []paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Currency: currency,
					Value:    amount,
				},
			},
		},
		ApplicationContext: &paypal.ApplicationContext{
			ReturnURL: returnURL,
			CancelURL: cancelURL,
		},
	}

	order, err := p.Client.CreateOrder(context.Background(), orderRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}

	return order, nil
}

// CapturePayment captures a PayPal order after user approval
func (p *PayPalProvider) CapturePayment(orderID string) (*paypal.CaptureOrderResponse, error) {
	captureResponse, err := p.Client.CaptureOrder(context.Background(), orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to capture PayPal payment: %w", err)
	}

	return captureResponse, nil
}

// RefundPayment issues a refund for a PayPal transaction
func (p *PayPalProvider) RefundPayment(captureID string, amount string, currency string) (*paypal.RefundResponse, error) {
	refundRequest := paypal.CaptureRefundRequest{
		Amount: &paypal.PurchaseUnitAmount{
			Currency: currency,
			Value:    amount,
		},
	}

	refundResponse, err := p.Client.RefundCapture(context.Background(), captureID, refundRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to process PayPal refund: %w", err)
	}

	return refundResponse, nil
}

// CreatePayout sends a payment to a recipient
func (p *PayPalProvider) CreatePayout(email string, amount string, currency string) (*paypal.PayoutBatch, error) {
	payout := &paypal.PayoutRequest{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: "You have received a payout!",
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      email,
				Amount: &paypal.PayoutAmount{
					Value:    amount,
					Currency: currency,
				},
			},
		},
	}

	payoutResponse, err := p.Client.CreatePayout(context.Background(), payout)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal payout: %w", err)
	}

	return payoutResponse, nil
}
