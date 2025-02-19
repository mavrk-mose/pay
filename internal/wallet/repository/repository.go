package repository

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
)

type WalletRepo interface {
	GetBalance(ctx context.Context, userID string) (float64, error)
	UpdateBalance(ctx context.Context, userID string, amount float64) error
	Create(ctx *gin.Context, wallet Wallet) error
	GetByID(ctx *gin.Context, userID string) (Wallet, error)
	CreateTransfer(ctx *gin.Context, transfer *TransferRequest) error
	UpdateTransferStatus(ctx *gin.Context, externalRef, status string) error
	Withdraw(ctx *gin.Context, walletID int64, amount float64, currency string) (string, error)
	Debit(ctx context.Context, walletID uuid.UUID, amount float64) error
	Credit(ctx context.Context, walletID uuid.UUID, amount float64) error
}

type walletRepo struct {
	DB *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) WalletRepo {
	return &walletRepo{DB: db}
}

func (r *walletRepo) GetBalance(ctx context.Context, userID string) (float64, error) {
	var balance float64
	err := r.DB.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error fetching wallet balance: %v", err)
	}
	return balance, nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID)
	return err
}

func (r *walletRepo) Create(ctx *gin.Context, wallet Wallet) error {
	_, err := r.DB.NamedExecContext(ctx, `INSERT INTO wallets (id, customer_id, balance, currency) 
		VALUES (:id, :customer_id, :balance, :currency)`, wallet)
	return err
}

func (r *walletRepo) GetByID(ctx *gin.Context, userID string) (Wallet, error) {
	wallet := &Wallet{}
	err := r.DB.GetContext(ctx, wallet, "SELECT * FROM wallets WHERE user_id = $1", userID)
	if err != nil {
		return Wallet{}, err
	}
	return *wallet, nil
}

func (r *walletRepo) CreateTransfer(ctx *gin.Context, transfer *TransferRequest) error {
	// Begin a transaction
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Debit the from_wallet
	debitQuery := "UPDATE wallets SET balance = balance - :amount WHERE id = :from_wallet_id"
	_, err = tx.NamedExecContext(ctx, debitQuery, transfer)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to debit from_wallet: %v", err)
	}

	// Credit the to_wallet
	creditQuery := "UPDATE wallets SET balance = balance + :amount WHERE id = :to_wallet_id"
	_, err = tx.NamedExecContext(ctx, creditQuery, transfer)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit to_wallet: %v", err)
	}

	// Insert the transfer record
	insertQuery := `
		INSERT INTO transfers (from_wallet_id, to_wallet_id, amount, currency, status, external_ref) 
		VALUES (:from_wallet_id, :to_wallet_id, :amount, :currency, :status, :external_ref)
	`
	_, err = tx.NamedExecContext(ctx, insertQuery, transfer)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create transfer record: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *walletRepo) UpdateTransferStatus(ctx *gin.Context, externalRef, status string) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE transfers SET status = $1 WHERE external_ref = $2", status, externalRef)
	return err
}

func (r *walletRepo) Withdraw(ctx *gin.Context, walletID int64, amount float64, currency string) (string, error) {
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

	_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, walletID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

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
	_, err = tx.ExecContext(ctx, creditQuery, transactionID, walletID, "CREDIT", amount, currency)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return transactionID, nil
}

// Debit subtracts funds from a wallet
func (r *walletRepo) Debit(ctx context.Context, walletID uuid.UUID, amount float64) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, walletID)
	return err
}

// Credit adds funds to a wallet
func (r *walletRepo) Credit(ctx context.Context, walletID uuid.UUID, amount float64) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
	return err
}
