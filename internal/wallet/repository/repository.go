package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
)

type WalletRepo struct {
	DB *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
	return &WalletRepo{DB: db}
}

func (r *WalletRepo) GetBalance(userID string) (float64, error) {
	var balance float64
	err := r.DB.QueryRow("SELECT balance FROM wallets WHERE user_id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error fetching wallet balance: %v", err)
	}
	return balance, nil
}

func (r *WalletRepo) UpdateBalance(userID string, amount float64) error {
	_, err := r.DB.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID)
	return err
}

func (r *WalletRepo) Create(ctx *gin.Context, wallet Wallet) error {
	_, err := r.DB.NamedExecContext(ctx, `INSERT INTO wallets (id, customer_id, balance, currency) 
		VALUES (:id, :customer_id, :balance, :currency)`, wallet)
	return err
}

func (r *WalletRepo) GetByID(ctx *gin.Context, userID string) (Wallet, error) {
	wallet := &Wallet{}
	err := r.DB.GetContext(ctx, wallet, "SELECT * FROM wallets WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *WalletRepo) UpdateWalletBalance(ctx *gin.Context, walletID uuid.UUID, amount float64) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
	return err
}

func (r *WalletRepo) CreateTransfer(ctx *gin.Context, transfer *TransferRequest) error {
	_, err := r.DB.NamedExecContext(ctx, `INSERT INTO transfers (from_wallet_id, to_wallet_id, amount, currency, status, external_ref) 
		VALUES (:from_wallet_id, :to_wallet_id, :amount, :currency, :status, :external_ref)`, transfer)
	return err
}

func (r *WalletRepo) UpdateTransferStatus(ctx *gin.Context, externalRef, status string) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE transfers SET status = $1 WHERE external_ref = $2", status, externalRef)
	return err
}

func (r *WalletRepo) Withdraw(ctx *gin.Context, walletID int64, amount float64, currency string) (string, error) {
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
