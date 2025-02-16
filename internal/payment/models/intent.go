package models

import "github.com/google/uuid"

type PaymentIntent struct {
	ID                   string               `json:"id"`
	Object               string               `json:"object"`
	Amount               float64              `json:"amount"`
	AmountCaptured       float64              `json:"amount_captured"`
	AmountRefunded       float64              `json:"amount_refunded"`
	AmountDetails        AmountDetails        `json:"amount_details"`
	AmountReceived       int64                `json:"amount_received"`
	Application          *string              `json:"application,omitempty"`
	ApplicationFee       *string              `json:"application_fee,omitempty"`
	ApplicationFeeAmount *int64               `json:"application_fee_amount,omitempty"`
	BalanceTransaction   string               `json:"balance_transaction"`
	BillingDetails       BillingDetails       `json:"billing_details"`
	Captured             bool                 `json:"captured"`
	Created              int64                `json:"created"`
	Currency             string               `json:"currency"`
	Customer             uuid.UUID            `json:"customer"`
	Description          string               `json:"description"`
	Disputed             bool                 `json:"disputed"`
	FailureCode          *string              `json:"failure_code,omitempty"`
	FailureMessage       *string              `json:"failure_message,omitempty"`
	Invoice              string               `json:"invoice"`
	Livemode             bool                 `json:"livemode"`
	Metadata             map[string]string    `json:"metadata"`
	PaymentIntent        string               `json:"payment_intent"`
	PaymentMethod        string               `json:"payment_method"`
	PaymentMethodOptions PaymentMethodOptions `json:"payment_method_options"`
	PaymentMethodTypes   []string             `json:"payment_method_types"`
	PaymentMethodDetails PaymentMethodDetails `json:"payment_method_details"`
	ReceiptEmail         string               `json:"receipt_email"`
	ReceiptNumber        string               `json:"receipt_number"`
	ReceiptURL           string               `json:"receipt_url"`
	Refunded             bool                 `json:"refunded"`
	Refunds              Refunds              `json:"refunds"`
	Status               string               `json:"status"`
}

type AmountDetails struct {
	Tip map[string]interface{} `json:"tip"`
}

// BillingDetails holds customer billing information.
type BillingDetails struct {
	Address Address `json:"address"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
}

// Address represents a postal address.
type Address struct {
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	PostalCode string  `json:"postal_code"`
	State      string  `json:"state"`
}

// Card represents card-specific details.
type Card struct {
	Brand        string  `json:"brand"`
	Checks       Checks  `json:"checks"`
	Country      string  `json:"country"`
	ExpMonth     int     `json:"exp_month"`
	ExpYear      int     `json:"exp_year"`
	Fingerprint  string  `json:"fingerprint"`
	Funding      string  `json:"funding"`
	Installments *string `json:"installments,omitempty"`
	Last4        string  `json:"last4"`
	Network      string  `json:"network"`
	ThreeDSecure *string `json:"three_d_secure,omitempty"`
	Wallet       *string `json:"wallet,omitempty"`
}

// PaymentMethodDetails holds information about the payment method.
type PaymentMethodDetails struct {
	Card Card   `json:"card"`
	Type string `json:"type"`
}

// Checks holds verification check results.
type Checks struct {
	AddressLine1Check      string `json:"address_line1_check"`
	AddressPostalCodeCheck string `json:"address_postal_code_check"`
	CVCCheck               string `json:"cvc_check"`
}

// Refunds represents a list of refunds.
type Refunds struct {
	Object     string        `json:"object"`
	Data       []interface{} `json:"data"` // Use []Refund if you want to define a Refund struct.
	HasMore    bool          `json:"has_more"`
	TotalCount int           `json:"total_count"`
	URL        string        `json:"url"`
}

type PaymentMethodOptions struct {
	Installments        interface{} `json:"installments"`
	MandateOptions      interface{} `json:"mandate_options"`
	Network             interface{} `json:"network"`
	RequestThreeDSecure string      `json:"request_three_d_secure"`
}

type PaymentMethodOptionsLink struct {
	PersistentToken interface{} `json:"persistent_token"`
}

// Stripe payment intent
// {
// 	"id": "pi_3MtwBwLkdIwHu7ix28a3tqPa",
// 	"object": "payment_intent",
// 	"amount": 2000,
// 	"amount_capturable": 0,
// 	"amount_details": {
// 	  "tip": {}
// 	},
// 	"amount_received": 0,
// 	"application": null,
// 	"application_fee_amount": null,
// 	"automatic_payment_methods": {
// 	  "enabled": true
// 	},
// 	"canceled_at": null,
// 	"cancellation_reason": null,
// 	"capture_method": "automatic",
// 	"client_secret": "pi_3MtwBwLkdIwHu7ix28a3tqPa_secret_YrKJUKribcBjcG8HVhfZluoGH",
// 	"confirmation_method": "automatic",
// 	"created": 1680800504,
// 	"currency": "usd",
// 	"customer": null,
// 	"description": null,
// 	"invoice": null,
// 	"last_payment_error": null,
// 	"latest_charge": null,
// 	"livemode": false,
// 	"metadata": {},
// 	"next_action": null,
// 	"on_behalf_of": null,
// 	"payment_method": null,
// 	"payment_method_options": {
// 	  "card": {
// 		"installments": null,
// 		"mandate_options": null,
// 		"network": null,
// 		"request_three_d_secure": "automatic"
// 	  },
// 	  "link": {
// 		"persistent_token": null
// 	  }
// 	},
// 	"payment_method_types": [
// 	  "card",
// 	  "link"
// 	],
// 	"processing": null,
// 	"receipt_email": null,
// 	"review": null,
// 	"setup_future_usage": null,
// 	"shipping": null,
// 	"source": null,
// 	"statement_descriptor": null,
// 	"statement_descriptor_suffix": null,
// 	"status": "requires_payment_method",
// 	"transfer_data": null,
// 	"transfer_group": null
//   }
