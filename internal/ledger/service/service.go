package ledger

import (
	. "github.com/mavrk-mose/pay/internal/fraud/models"
)

// Ledger module (immutable transactions)
type LedgerService interface {
	RecordTransaction(transaction Transaction) error
	GetTransactionByID(transactionID string) (Transaction, error)
}
