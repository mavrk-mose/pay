package services

import (
	. "github.com/mavrk-mose/pay/internal/model"
)

type LedgerService interface {
	RecordTransaction(transaction Transaction) error
	GetTransactionByID(transactionID string) (Transaction, error)
	ReconcileTransactions(settlementFile SettlementFile) (ReconciliationResult, error)
}
