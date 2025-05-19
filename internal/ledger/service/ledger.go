package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	models "github.com/mavrk-mose/pay/internal/ledger/models"
)

//go:generate mockery --name=LedgerService --output=./mocks --filename=ledger.go --with-expecter
type LedgerService interface {
	RecordTransaction(ctx *gin.Context, txn models.Transaction) (string, error)
	CreateTransactionWithEntries(ctx *gin.Context, txn *sqlx.Tx, entries []models.Transaction) error
	UpdateTransactionStatus(ctx *gin.Context, externalRef string, status models.TransactionStatus) error
	FetchTransactionsWithChecksum(db *sqlx.DB, date, provider string) (map[string]string, error)
	GetTransactionByID(transactionID string) (models.Transaction, error)
}
