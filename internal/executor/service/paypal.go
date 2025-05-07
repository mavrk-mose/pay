package service

import (
	"encoding/json"
	"fmt"

	. "github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/mavrk-mose/pay/pkg/utils"
	paypalsdk "github.com/netlify/PayPal-Go-SDK"
)

type PayPalProvider struct {
	Client 		 *paypalsdk.Client
	httpClient   utils.GenericHttpClient
	logger       utils.Logger
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

func (p *PayPalProvider) ExecutePayment(order PaymentOrder) (any, error) { 
	panic("unimplemented")
}

func (p *PayPalProvider) CreateOrder(amount string, currency string, returnURL string, cancelURL string) (*paypalsdk.Order, error) {
	orderReq := paypalsdk.CreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []paypalsdk.PurchaseUnitRequest{
			{
				Amount: &paypalsdk.PurchaseUnitAmount{
					Currency: currency,
					Value:    amount,
				},
			},
		},
		ApplicationContext: &paypalsdk.ApplicationContext{
			ReturnURL: returnURL,
			CancelURL: cancelURL,
		},
	}

	order, err := p.Client.CreateOrder(orderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}

	return order, nil
}

// CapturePayment captures a PayPal order after user approval
func (p *PayPalProvider) CapturePayment(orderID string) (*paypalsdk.Capture, error) {
	captureResponse, err := p.Client.CaptureOrder(orderID, nil, true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to capture PayPal payment: %w", err)
	}

	return captureResponse, nil
}

// RefundPayment issues a refund for a PayPal transaction
func (p *PayPalProvider) RefundPayment(captureID string, amount string, currency string) (*paypalsdk.Refund, error) {
	// refundRequest := p.Client.RefundRequest{
	// 	Amount: &p.Client.PurchaseUnitAmount{
	// 		Currency: currency,
	// 		Value:    amount,
	// 	},
	// }

	// refundResponse, err := p.Client.RefundCapture(context.Background(), captureID, refundRequest)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to process PayPal refund: %w", err)
	// }

	// return refundResponse, nil

	panic("unimplemented")
}

// CreatePayout sends a payment to a recipient
func (p *PayPalProvider) CreatePayout(email string, amount string, currency string) (*paypalsdk.Payout, error) { // supposed to be payout in batches
	// payout := &p.Client.PayoutRequest{ // supposed to be payout in batches
	// 	SenderBatchHeader: &p.Client.SenderBatchHeader{
	// 		EmailSubject: "You have received a payout!",
	// 	},
	// 	Items: []p.Client.PayoutItem{
	// 		{
	// 			RecipientType: "EMAIL",
	// 			Receiver:      email,
	// 			Amount: &p.Client.PayoutAmount{
	// 				Value:    amount,
	// 				Currency: currency,
	// 			},
	// 		},
	// 	},
	// }

	// payoutResponse, err := p.Client.CreatePayout(context.Background(), payout)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create PayPal payout: %w", err)
	// }

	// return payoutResponse, nil

	panic("unimplemented")
}

// PayPalWebhookVerifier verifies webhooks from PayPal
func (p *PayPalProvider) VerifyWebhook(headers map[string]string, body []byte, webhookID string) (bool, error) {
	req := map[string]any{
		"auth_algo":         headers["PAYPAL-AUTH-ALGO"],
		"cert_url":          headers["PAYPAL-CERT-URL"],
		"transmission_id":   headers["PAYPAL-TRANSMISSION-ID"],
		"transmission_sig":  headers["PAYPAL-TRANSMISSION-SIG"],
		"transmission_time": headers["PAYPAL-TRANSMISSION-TIME"],
		"webhook_id":        webhookID,
		"webhook_event":     json.RawMessage(body),
	}

	accessTokenResp, err := p.Client.GetAccessToken()
	if err != nil {
		return false, fmt.Errorf("failed to get access token: %w", err)
	}

	// TODO: get the paypal URL through the config
	response, err := p.httpClient.Post("https://api-m.sandbox.paypal.com/v1/notifications/verify-webhook-signature", req, map[string]string{
		"Authorization": "Bearer " + accessTokenResp.Token,
	})
	if err != nil {
		return false, err
	}

	var res map[string]string
	if err := json.Unmarshal(*response, &res); err != nil {
		return false, err
	}

	return res["verification_status"] == "SUCCESS", nil
}
