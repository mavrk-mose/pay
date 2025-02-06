package model

type LedgerTransaction struct {
	TransactionID string
	DebitAccount  string
	CreditAccount string
	Amount        float64
	Timestamp     string
}

type WalletTransaction struct {
	TransactionID string
	UserID        string
	Amount        float64
	Timestamp     string
}
