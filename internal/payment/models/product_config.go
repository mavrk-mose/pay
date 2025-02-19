package models

import "github.com/google/uuid"

type ProductConfiguration struct {
	ID            uuid.UUID `db:"id"`
	ProductName   string    `db:"product_name"`
	FeePercentage float64   `db:"fee_percentage"`
	FeeWalletID   uuid.UUID `db:"fee_wallet_id"`
	// any other relevant fields
}
