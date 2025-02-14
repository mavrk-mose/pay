package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	Active     Status = "active"
	Terminated Status = "terminated"
)

type Wallet struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`
	Balance    float64   `db:"balance" json:"balance"`
	Status     Status    `db:"status" json:"status,omitempty"`
	Currency   string    `db:"currency" json:"currency"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
