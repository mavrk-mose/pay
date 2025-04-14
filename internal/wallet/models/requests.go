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

type WalletTransactionRequest struct {
	WalletID    string  `json:"walletID" binding:"required"`    // the target wallet
	Amount      float64 `json:"amount" binding:"required,gt=0"` // must be greater than 0
	Currency    string  `json:"currency" binding:"required"`    // e.g., "USD"
	Description string  `json:"description,omitempty"`          // optional description
	ReferenceID string  `json:"referenceID,omitempty"`          // optional for idempotency tracking
}