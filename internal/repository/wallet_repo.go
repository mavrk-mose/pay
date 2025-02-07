package repository

import (
    "context"
    "database/sql"
    "errors"
    "github.com/jmoiron/sqlx"
	"github.com/google/uuid"
    "github.com/mavrk-mose/pay/internal/model"
)

type WalletRepo struct {
    DB *sqlx.DB
	GetByID(ctx context.Context, userID int64) (*models.Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID int64, amount float64) error
	CreateTransfer(ctx context.Context, transfer *models.Transfer) error
	UpdateTransferStatus(ctx context.Context, externalRef, status string) error
	Withdraw(ctx context.Context, walletID int64, amount float64, currency string) (string, error)
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
    return &WalletRepo{DB: db}
}

func (r *WalletRepo) GetByID(ctx context.Context, userID int64) (*models.Wallet, error) {
    wallet := &models.Wallet{}
    err := r.DB.GetContext(ctx, wallet, "SELECT * FROM wallets WHERE user_id = $1", userID)
    if err != nil {
        return nil, err
    }
    return wallet, nil
}

func (r *WalletRepo) UpdateWalletBalance(ctx context.Context, walletID int64, amount float64) error {
    _, err := r.DB.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
    return err
}

func (r *WalletRepo) CreateTransfer(ctx context.Context, transfer *models.Transfer) error {
    _, err := r.DB.NamedExecContext(ctx, `INSERT INTO transfers (from_wallet_id, to_wallet_id, amount, currency, status, external_ref) 
        VALUES (:from_wallet_id, :to_wallet_id, :amount, :currency, :status, :external_ref)`, transfer)
    return err
}

func (r *WalletRepo) UpdateTransferStatus(ctx context.Context, externalRef, status string) error {
    _, err := r.DB.ExecContext(ctx, "UPDATE transfers SET status = $1 WHERE external_ref = $2", status, externalRef)
    return err
}

func (r *WalletRepo) Withdraw(ctx context.Context, walletID int64, amount float64, currency string) (string, error) {
	
	transactionID := uuid.New().String()

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	var currentBalance float64
	err = tx.GetContext(ctx, &currentBalance, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	if currentBalance < amount {
		tx.Rollback()
		return "", fmt.Errorf("insufficient funds: available %v, required %v", currentBalance, amount)
	}

	// Update the customer's wallet: debit the withdrawal amount.
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, walletID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// Update the withdrawal account: credit the withdrawal amount.
	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, r.WithdrawalAccount)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	// In this record, the customer's wallet is debited, and the withdrawal account is credited.
	debitQuery := `
		INSERT INTO transaction (transaction_id, wallet_id, entry_type, amount, currency)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.ExecContext(ctx, debitQuery, transactionID, walletID, "DEBIT", amount, currency)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	creditQuery := `
		INSERT INTO transaction (transaction_id, wallet_id, entry_type, amount, currency)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, creditQuery, transactionID, r.WithdrawalAccount, "CREDIT", amount, currency)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	var txnID int64
	err = tx.QueryRowxContext(ctx, transactionInsert,
		transactionID,
		model.TransactionWithdrawal,
		model.TransactionPending,
		currency,
		walletID,
		amount,
		r.WithdrawalAccount,
		amount,
	).Scan(&txnID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return transactionID, nil
}
