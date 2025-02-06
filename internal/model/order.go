package model

import (
	"time"
)

type PaymentStatus string

func (p PaymentStatus) UnmarshalText(b []byte) any {
	panic("unimplemented")
}

const (
	StatusPending   PaymentStatus = "pending"
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
	StatusRefunded  PaymentStatus = "refunded"
)

type PaymentOrder struct {
	ID              string        `json:"id"`               // Unique identifier for the payment
	Amount          float64       `json:"amount"`           // Payment amount
	Currency        string        `json:"currency"`         // Payment currency (e.g., USD, EUR)
	PayerID         string        `json:"payer_id"`         // Identifier for the payer
	PayeeID         string        `json:"payee_id"`         // Identifier for the payee
	PaymentMethod   string        `json:"payment_method"`   // Method of payment (e.g., credit_card, bank_transfer)
	Status          PaymentStatus `json:"status"`           // Status of the payment (e.g., pending, completed, failed)
	ReferenceNumber string        `json:"reference_number"` // Reference number for external systems
	CreatedAt       time.Time     `json:"created_at"`       // Timestamp when the payment was created
	UpdatedAt       time.Time     `json:"updated_at"`       // Timestamp when the payment was last updated
	Description     string        `json:"description"`      // Optional description or memo for the payment
}

