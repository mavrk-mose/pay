package ports

import . "github.com/mavrk-mose/pay/ledger/internal/model"

type LedgerService interface {
	RecordTransaction(transaction LedgerTransaction) error
	GetTransactionByID(transactionID string) (LedgerTransaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}
