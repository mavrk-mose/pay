package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/internal/model"
)

type TransactionRepo interface {
	CreateTransactionRecord(ctx context.Context, rec *model.Transaction) error
	UpdateTransactionStatus(ctx context.Context, externalRef string, status model.TransactionStatus) error
}

type pgTransactionRepo struct {
	DB *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) TransactionRepo {
	return &pgTransactionRepo{DB: db}
}

func (r *pgTransactionRepo) CreateTransactionRecord(ctx context.Context, txn *model.Transaction) error {
	query := `
		INSERT INTO transaction (
			external_ref, type, status, details, currency, 
			debit_wallet_id, debit_amount, credit_wallet_id, 
			credit_amount
			) VALUES (
				:external_ref, :type, :status, :details, 
				:currency, :debit_wallet_id, :debit_amount, 
				:credit_wallet_id, :credit_amount
			)
		RETURNING id, created_at, updated_at
	`
	return r.DB.QueryRowxContext(ctx, query, txn).
		Scan(&txn.ID, &txn.CreatedAt, &txn.UpdatedAt)
}

func (r *pgTransactionRepo) UpdateTransactionStatus(ctx context.Context, externalRef string, status model.TransactionStatus) error {
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
