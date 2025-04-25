package service

import (
	"github.com/gin-gonic/gin"
	models "github.com/mavrk-mose/pay/internal/ledger/models"
)

// LedgerService Ledger module (immutable transactions)
type LedgerService interface {
	RecordTransaction(ctx *gin.Context, txn models.Transaction) error
	GetTransactionByID(transactionID string) (models.Transaction, error)
}
