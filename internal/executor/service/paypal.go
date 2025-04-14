package service

import (
	"context"
	"fmt"
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/mavrk-mose/pay/pkg/utils"
	paypalsdk "github.com/netlify/PayPal-Go-SDK"
)

type PayPalProvider struct {
	Client *paypalsdk.Client
	logger utils.Logger
}

func NewPayPalProvider(clientID, secret string, isSandbox bool) (*PayPalProvider, error) {
	var apiBase string
	if isSandbox {
		apiBase = paypalsdk.APIBaseSandBox
	} else {
		apiBase = paypalsdk.APIBaseLive
	}

	client, err := paypalsdk.NewClient(clientID, secret, apiBase)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal client: %w", err)
	}

	_, err = client.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve PayPal access token: %w", err)
	}

	return &PayPalProvider{Client: client}, nil
}

func (p *PayPalProvider) CreateOrder(amount string, currency string, returnURL string, cancelURL string) (*paypalsdk.Order, error) {
	orderIntent := p.Client.IntentCapture()
	orderRequest := p.Client.OrderRequest{
		Intent: orderIntent,
		PurchaseUnits: []p.Client.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Currency: currency,
					Value:    amount,
				},
			},
		},
		ApplicationContext: &p.Client{
			ReturnURL: returnURL,
			CancelURL: cancelURL,
		},
	}

	order, err := p.Client.CreatePayment(context.Background(), orderRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}

	return order, nil
}

// CapturePayment captures a PayPal order after user approval
func (p *PayPalProvider) CapturePayment(orderID string) (*paypalsdk.CaptureOrderResponse, error) {
	captureResponse, err := p.Client.CaptureOrder(context.Background(), orderID, paypalsdk.CaptureOrderRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to capture PayPal payment: %w", err)
	}

	return captureResponse, nil
}

// RefundPayment issues a refund for a PayPal transaction
func (p *PayPalProvider) RefundPayment(captureID string, amount string, currency string) (*paypalsdk.RefundResponse, error) {
	refundRequest := p.Client.RefundRequest{
		Amount: &p.Client.PurchaseUnitAmount{
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
func (p *PayPalProvider) CreatePayout(email string, amount string, currency string) (*paypalsdk.Payout, error) { // supposed to be payout in batches
	payout := &p.Client.PayoutRequest{ // supposed to be payout in batches
		SenderBatchHeader: &p.Client.SenderBatchHeader{
			EmailSubject: "You have received a payout!",
		},
		Items: []p.Client.PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      email,
				Amount: &p.Client.PayoutAmount{
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
