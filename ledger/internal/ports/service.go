package ports

type LedgerService interface {
    RecordTransaction(transaction LedgerTransaction) error
    GetTransactionByID(transactionID string) (LedgerTransaction, error)
    ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}