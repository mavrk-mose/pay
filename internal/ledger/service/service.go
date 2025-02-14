package ledger

// Ledger module (immutable transactions)
type LedgerService interface {
	RecordTransaction(transaction Transaction) error
	GetTransactionByID(transactionID string) (Transaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}

