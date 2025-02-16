package ledger

import (
	"github.com/gin-gonic/gin"
	. "github.com/mavrk-mose/pay/internal/fraud/models"
)

// Ledger module (immutable transactions)
type LedgerService interface {
	RecordTransaction(ctx *gin.Context, txn Transaction) error
	GetTransactionByID(transactionID string) (Transaction, error)
}
