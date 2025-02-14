package wallet

type CreateWalletRequest struct {
	CustomerID     string  `json:"customer_id"`
	InitialBalance float64 `json:"initial_balance"`
	Currency       string  `json:"currency"`
}

type TransferRequest struct {
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
}

type WithdrawalRequest struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type DepositRequest struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
