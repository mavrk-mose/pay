package model

type LedgerTransaction struct {
	TransactionID string
	DebitAccount  string
	CreditAccount string
	Amount        float64
	Timestamp     string
}
