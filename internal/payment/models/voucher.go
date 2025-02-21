package models

import(
	"time"
)

type Voucher struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Code      string    `db:"code" json:"code"`
	Amount    float64   `db:"amount" json:"amount"`
	Currency  string    `db:"currency" json:"currency"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Redeemed  bool      `db:"redeemed" json:"redeemed"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}