package ledger

import (
	"fmt"
	. "github.com/mavrk-mose/pay/internal/fraud/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

func (r *Repo) RecordTransaction(ctx *gin.Context, payerWalletID, payeeWalletID int64, amount float64, currency string) (string, error) {
	txn, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	transactionID := uuid.New()

	// Create ledger entries: one debit and one credit
	debitEntry := Transaction{
		ID:            transactionID,
		ExternalRef:   uuid.New().String(),
		Type:          TransactionTransfer,
		Status:        TransactionPending,
		Currency:      currency,
		DebitWalletID: payerWalletID,
		DebitAmount:   amount,
		EntryType:     Debit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	creditEntry := Transaction{
		ID:             transactionID,
		ExternalRef:    uuid.New().String(),
		Type:           TransactionTransfer,
		Status:         TransactionPending,
		Currency:       currency,
		CreditWalletID: payeeWalletID,
		CreditAmount:   amount,
		EntryType:      Credit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	entries := []Transaction{debitEntry, creditEntry}

	if err := r.CreateTransactionWithEntries(ctx, txn, entries); err != nil {
		err := txn.Rollback()
		if err != nil {
			return "", err
		}
		return "", err
	}

	if err := txn.Commit(); err != nil {
		return "", err
	}

	return transactionID.String(), nil
}

func (r *Repo) CreateTransactionWithEntries(ctx *gin.Context, txn *sqlx.Tx, entries []Transaction) error {
	for _, entry := range entries {
		_, err := txn.NamedExecContext(ctx, `
			INSERT INTO transactions (
				id, external_ref, type, status, details, currency, 
				debit_wallet_id, debit_amount, entry_type, 
				credit_wallet_id, credit_amount, created_at, updated_at
			) VALUES (
				:id, :external_ref, :type, :status, :details, :currency, 
				:debit_wallet_id, :debit_amount, :entry_type, 
				:credit_wallet_id, :credit_amount, :created_at, :updated_at
			)`, entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) UpdateTransactionStatus(ctx *gin.Context, externalRef string, status TransactionStatus) error {
	query := `
		UPDATE transaction
		SET status = $1, updated_at = NOW() 
		WHERE external_ref = $2
	`
	res, err := r.DB.ExecContext(ctx, query, status, externalRef)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return fmt.Errorf("no transaction found with external_ref %s", externalRef)
	}
	return nil
}
