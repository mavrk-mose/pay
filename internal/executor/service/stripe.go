package service

import (
	. "github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/checkout/session"
)

type StripeProvider struct {
    SecretKey string // load from yaml config
}

func NewStripeProvider(secretKey string) *StripeProvider {
	return &StripeProvider{SecretKey: secretKey}
}

func (s *StripeProvider) ExecutePayment(order PaymentOrder) (any, error) {
	// stripe.Key = s.SecretKey

	// params := &stripe.PaymentIntentParams{
	// 	Amount:   stripe.Int64(order.Amount), // amount in cents
	// 	Currency: stripe.String(order.Currency),
	// 	Confirm:  stripe.Bool(true),
	// 	PaymentMethodTypes: stripe.StringSlice([]string{
	// 		"card",
	// 	}),
	// }

	// intent, err := paymentintent.New(params)
	// if err != nil {
	// 	return PaymentResult{}, fmt.Errorf("failed to create payment intent: %w", err)
	// }

	// return PaymentResult{
	// 	Success:       true,
	// 	Message:       "Payment processed successfully",
	// 	TransactionID: intent.ID,
	// }, nil
    panic("unimplemented")
}

func (s *StripeProvider) CreateCustomer() (*stripe.Customer, error) {
    stripe.Key = s.SecretKey
    params := &stripe.CustomerParams{
        Name:  stripe.String("John Dey"),
        Email: stripe.String("abcd123@gmail.com"),
        Address: &stripe.AddressParams{
        Line1:      stripe.String("1234 Elm St"),
        City:       stripe.String("Smalltown"),
        PostalCode: stripe.String("12345"),
        State:      stripe.String("CA"),
        Country:    stripe.String("US"),
        },
    }
    result, err := customer.New(params)
    return result, err
}

func (s *StripeProvider) CreateCheckoutSession() (*stripe.CheckoutSession, error) {
    stripe.Key = s.SecretKey
    params := &stripe.CheckoutSessionParams{
        SuccessURL: stripe.String("https://example.com/success"), 
        Mode:       stripe.String("setup"),                      
        Customer:   stripe.String("cus_HKtmyFxyxPZQDm"),  
        PaymentMethodTypes: stripe.StringSlice([]string{
        "card",
        }),
    }
    result, err := session.New(params)
    return result, err
}

// {
//     "data": {
//         "sessionId": "cs_test_c1vTIqOhhs9db05IFiZuanbTcHgboHu2aOLbpjsbVagkY08h4muUJ07Nmo",
//         "successUrl": "http://localhost:8088/session/successfull",
//         "saveInfoUrl": "https://checkout.stripe.com/c/pay/cs_test_c1vTIqOhhs9db05IFiZuanbTcHgboHu2aOLbpjsbVagkY08h4muUJ07Nmo#fidkdWxOYHwnPyd1blpxYHZxWjA0SnR3ZE5WT2JUfEN2ZmQ1RkJ2NGdddn10QUJNclUza3FfY2xffHxsVWdBT1ZiYWxLXUEzPTx%2FYWI1PXZqdXMyMlZsZ0liQ3ZoNHVHZkkktAZ013c0t1cm5NNTVhRGl2RHBDUycpJ2N3amhWYHdzYHcnP3F3cGApJ2lkfGpwcVF8dWAnPyd2bGtiaWBaZmppcGhrJyknYGtkZ2lgVWlkZmBtamlhYHd2Jz9xd3BgeCUl"
//     },
//     "message": "Customer checkout done successfully !!"
// }

func (s *StripeProvider) CreatePaymentIntent() (*stripe.PaymentIntent, error) {
    stripe.Key = s.SecretKey 
    params := &stripe.PaymentIntentParams{
        Amount:      stripe.Int64(1000000), // Assuming `Balance` from the request is the amount in cents
        Currency:    stripe.String("usd"), // Currency such as "usd"
        Confirm:     stripe.Bool(true),
        PaymentMethod: stripe.String("payment_method_id after saving card"),
        // PaymentMethodTypes: stripe.StringSlice(["card"]),
        // Customer:    stripe.String(customer_id), 
        OffSession:  stripe.Bool(true),
        ConfirmationMethod: stripe.String("automatic"), 
        Description: stripe.String("Example payment intent for invoice"), // Description of the transaction
        Shipping: &stripe.ShippingDetailsParams{
            Name: stripe.String("John Doe"),
            Address: &stripe.AddressParams{
              Line1: stripe.String("1234 Main Street"),
              PostalCode: stripe.String("94105"),
              City: stripe.String("San Francisco"),
              State: stripe.String("CA"),
              Country: stripe.String("US"),
            },
        },
    }

    result, err := paymentintent.New(params)
    return result, err
}

func (s *StripeProvider) CreateCharge() (*stripe.Charge, error) {
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(2000),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// Desc:     stripe.String("Test Charge"),
	}
	params.SetSource("tok_visa") // use a test card token provided by Stripe
	ch, err := charge.New(params)
	return ch, err
}

