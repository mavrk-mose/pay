package repository

import (
	"context"
	"github.com/mavrk-mose/pay/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LedgerRepo struct {
	DB *sqlx.DB
}

func (r *LedgerRepo) RecordTransaction(ctx *gin.Context, payerWalletID, payeeWalletID int64, amount float64, currency string) (string, error) {
	txn, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	transactionID := uuid.New()

	// Create ledger entries: one debit and one credit
	debitEntry := model.Transaction{
		ID:            transactionID,
		ExternalRef:   uuid.New().String(),
		Type:          model.TransactionTransfer,
		Status:        model.TransactionPending,
		Currency:      currency,
		DebitWalletID: payerWalletID,
		DebitAmount:   amount,
		EntryType:     model.Debit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	creditEntry := model.Transaction{
		ID:             transactionID,
		ExternalRef:    uuid.New().String(),
		Type:           model.TransactionTransfer,
		Status:         model.TransactionPending,
		Currency:       currency,
		CreditWalletID: payeeWalletID,
		CreditAmount:   amount,
		EntryType:      model.Credit,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	entries := []model.Transaction{debitEntry, creditEntry}

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

func (r *LedgerRepo) CreateTransactionWithEntries(ctx context.Context, txn *sqlx.Tx, entries []model.Transaction) error {
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
