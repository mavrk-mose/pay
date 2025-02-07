package services

type PaymentExecutorService interface {
	ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error)
	RecordPaymentOrder(order PaymentOrder) error
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

// NewDefaultPaymentExecutor initializes the default gateways (e.g., Stripe, PayPal, Adyen)
// and returns a PaymentExecutorService instance.
func NewDefaultPaymentExecutor() PaymentExecutorService {
	gateways := map[string]PaymentGateway{
		"stripe": &StripeGateway{},
		"paypal": &PaypalGateway{},
		"adyen":  &AdyenGateway{},
	}
	return NewPaymentExecutor(gateways)
}


// ExecutePayment routes the payment request to the appropriate gateway based on order.Gateway.
// TODO: should use payment intent here
func (pe *PaymentExecutor) ExecutePayment(order PaymentOrder) (PaymentExecutionResult, error) {
	gateway, exists := pe.gateways[order.Gateway]
	if !exists {
		return PaymentExecutionResult{}, errors.New("unsupported payment gateway: " + order.Gateway)
	}
	return gateway.ExecutePayment(order)
}