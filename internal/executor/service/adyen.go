package service

import (
	"github.com/mavrk-mose/pay/pkg/utils"
	"github.com/adyen/adyen-go-api-library/v7/src/checkout"
	"github.com/adyen/adyen-go-api-library/v7/src/common"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type AdyenProvider struct {
	Client *checkout.APIClient
	MerchantAccount string
	logger          utils.Logger
}

func NewAdyenProvider(apiKey, merchantAccount string, isLive bool) (*AdyenProvider, error) {
	// env := common.TestEnv
	// if isLive {
	// 	env = common.LiveEnv
	// }

	client := checkout.NewAPIClient( &common.Client{
		
    })

	return &AdyenProvider{Client: client, MerchantAccount: merchantAccount}, nil
}

// ExecutePayment processes a payment request using Adyen
func (a *AdyenProvider) ExecutePayment(order PaymentOrder) (any, error) {
	// request := checkout.PaymentRequest{
	// 	AccountInfo:      nil,
	// 	AdditionalAmount: nil,
	// 	AdditionalData:   nil,
	// 	Amount: checkout.Amount{
	// 		Currency: order.Currency,
	// 		Value:    int64(order.Amount * 100), // Convert to minor units
	// 	},
	// 	ApplicationInfo:           nil,
	// 	AuthenticationData:        nil,
	// 	BillingAddress:            nil,
	// 	BrowserInfo:               nil,
	// 	CaptureDelayHours:         nil,
	// 	Channel:                   nil,
	// 	CheckoutAttemptId:         nil,
	// 	Company:                   nil,
	// 	ConversionId:              nil,
	// 	CountryCode:               nil,
	// 	DateOfBirth:               nil,
	// 	DccQuote:                  nil,
	// 	DeliveryAddress:           nil,
	// 	DeliveryDate:              nil,
	// 	DeviceFingerprint:         nil,
	// 	EnableOneClick:            nil,
	// 	EnablePayOut:              nil,
	// 	EnableRecurring:           nil,
	// 	EntityType:                nil,
	// 	FraudOffset:               nil,
	// 	IndustryUsage:             nil,
	// 	Installments:              nil,
	// 	LineItems:                 nil,
	// 	LocalizedShopperStatement: nil,
	// 	Mandate:                   nil,
	// 	Mcc:                       nil,
	// 	MerchantAccount:           a.MerchantAccount,
	// 	MerchantOrderReference:    nil,
	// 	MerchantRiskIndicator:     nil,
	// 	Metadata:                  nil,
	// 	MpiData:                   nil,
	// 	Order:                     nil,
	// 	OrderReference:            nil,
	// 	Origin:                    nil,
	// 	PaymentMethod:				nil,
	// 	PlatformChargebackLogic:   nil,
	// 	RecurringExpiry:           nil,
	// 	RecurringFrequency:        nil,
	// 	RecurringProcessingModel:  nil,
	// 	RedirectFromIssuerMethod:  nil,
	// 	RedirectToIssuerMethod:    nil,
	// 	Reference:                 order.OrderID,
	// 	ReturnUrl:                 "",
	// 	RiskData:                  nil,
	// 	SessionValidity:           nil,
	// 	ShopperEmail:              nil,
	// 	ShopperIP:                 nil,
	// 	ShopperInteraction:        nil,
	// 	ShopperLocale:             nil,
	// 	ShopperName:               nil,
	// 	ShopperReference:          nil,
	// 	ShopperStatement:          nil,
	// 	SocialSecurityNumber:      nil,
	// 	Splits:                    nil,
	// 	Store:                     nil,
	// 	StorePaymentMethod:        nil,
	// 	TelephoneNumber:           nil,
	// 	ThreeDS2RequestData:       nil,
	// 	ThreeDSAuthenticationOnly: nil,
	// 	TrustedShopper:            nil,
	// };

	// response, httpResp, err := a.Client.PaymentsApi.Payments(context.Background(), request)
	// if err != nil {
	// 	return nil, fmt.Errorf("payment failed: %w", err)
	// }
	// defer httpResp.Body.Close()

	// return &response, nil
	panic("unimplemented")
}

// CapturePayment captures an authorized payment
func (a *AdyenProvider) CapturePayment(paymentID string, amount float64, currency string) (*checkout.ModificationsApi, error) {
	// request := checkout.PaymentCaptureRequest{
	// 	Amount: checkout.Amount{
	// 		Currency: currency,
	// 		Value:    int64(amount * 100),
	// 	},
	// 	MerchantAccount: a.MerchantAccount,
	// }

	// response, httpResp, err := a.Client.PaymentsApi.Payments(context.Background(), request)
	// if err != nil {
	// 	return nil, fmt.Errorf("capture failed: %w", err)
	// }
	// defer httpResp.Body.Close()

	// return &response, nil
	panic("unimplemented")
}

// RefundPayment issues a refund for a payment
func (a *AdyenProvider) RefundPayment(paymentID string, amount float64, currency string) (any, error) {
	// request := checkout.RefundRequest{
	// 	Amount: checkout.Amount{
	// 		Currency: currency,
	// 		Value:    int64(amount * 100),
	// 	},
	// 	MerchantAccount: a.MerchantAccount,
	// }

	panic("unimplemented")
}

// CancelPayment cancels an authorized payment
func (a *AdyenProvider) CancelPayment(paymentID string) (any, error) {
	// request := checkout.CancelOrderRequest{}

	// response, httpResp, err := a.Client.PaymentsApi.Payments(context.Background(), paymentID, &request)
	// if err != nil {
	// 	return nil, fmt.Errorf("cancel failed: %w", err)
	// }
	// defer httpResp.Body.Close()

	// return &response, nil
	panic("unimplemented")
}
