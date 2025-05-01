package service

import (
	"github.com/gin-gonic/gin"
	models "github.com/mavrk-mose/pay/internal/ledger/models"
)

//go:generate mockery --name=LedgerService --output=./mocks --filename=ledger.go --with-expecter
type LedgerService interface {
	RecordTransaction(ctx *gin.Context, txn models.Transaction) error
	GetTransactionByID(transactionID string) (models.Transaction, error)
}
