package model

type SettlementFile struct {
	FileID        string
	BankAccountID string
	Transactions  []LedgerTransaction
}
