package ports

import . "github.com/mavrk-mose/pay/reconcile/internal/model"

type ReconciliationService interface {
	ReconcileLedgerWithSettlement(ledger LedgerService, settlementFile SettlementFile) (ReconciliationResult, error)
}

type LedgerService interface {
	RecordTransaction(transaction LedgerTransaction) error
	GetTransactionByID(transactionID string) (LedgerTransaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}
