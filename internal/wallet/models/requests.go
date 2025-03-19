package wallet

import "github.com/google/uuid"

type CreateWalletRequest struct {
	CustomerID string `json:"customer_id"`
	Currency   string `json:"currency"`
}

type TransferRequest struct {
	FromWalletID uuid.UUID `json:"from_wallet_id"`
	ToWalletID   uuid.UUID `json:"to_wallet_id"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
}

type WithdrawalRequest struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
}

type DepositRequest struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
}
