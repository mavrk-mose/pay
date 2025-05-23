package repository

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/internal/wallet/service"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type walletRepo struct {
	DB     *sqlx.DB
	logger utils.Logger
}

func NewWalletService(db *sqlx.DB) service.WalletService {
	return &walletRepo{DB: db}
}

func (r *walletRepo) CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error) {
	r.logger.Infof("Creating wallet for user %s", wallet.UserId)
	query := `
		INSERT INTO wallets (user_id, balance, currency, created_at)
		VALUES ($1, 0.00, $2, NOW())
		RETURNING id, user_id, balance, currency, created_at
	`
	err := r.DB.QueryRowx(query, wallet.UserId, wallet.Currency).StructScan(&wallet)
	if err != nil {
		r.logger.Errorf("Failed to create wallet: %v", err)
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}
	return wallet, nil
}

func (r *walletRepo) DeleteWallet(c context.Context, walletID string) error {
	r.logger.Infof("Deleting wallet %s", walletID)
	query := `DELETE FROM wallets WHERE id = $1`
	_, err := r.DB.ExecContext(c, query, walletID)
	if err != nil {
		r.logger.Errorf("Failed to delete wallet: %v", err)
		return fmt.Errorf("failed to delete wallet: %v", err)
	}
	return nil
}

func (r *walletRepo) GetUserWallets(ctx context.Context, userID string) ([]Wallet, error) {
	r.logger.Infof("Fetching wallets for user %s", userID)
	var wallets []Wallet
	query := `SELECT id, user_id, balance, currency, created_at FROM wallets WHERE user_id = $1`
	err := r.DB.Select(&wallets, query, userID)
	if err != nil {
		r.logger.Errorf("Failed to fetch wallets: %v", err)
		return nil, fmt.Errorf("failed to fetch wallets: %v", err)
	}
	return wallets, nil
}

func (r *walletRepo) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	r.logger.Infof("Fetching wallet balance for user %s", userID)
	var balance float64
	err := r.DB.QueryRowContext(ctx, "SELECT balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		r.logger.Errorf("Failed to fetch wallet balance: %v", err)
		return 0, fmt.Errorf("error fetching wallet balance: %v", err)
	}
	return balance, nil
}

func (r *walletRepo) GetByID(ctx context.Context, walletID uuid.UUID) (Wallet, error) {
	r.logger.Infof("Fetching wallet details for wallet %s", walletID)
	var wallet Wallet
	err := r.DB.GetContext(ctx, &wallet, "SELECT * FROM wallets WHERE id = $1", walletID)
	if err != nil {
		r.logger.Errorf("Failed to fetch wallet: %v", err)
		return Wallet{}, err
	}
	return wallet, nil
}

func (r *walletRepo) CreateTransfer(ctx context.Context, transfer *TransferRequest) error {
	r.logger.Infof("Initiating transfer for wallet %s", transfer.FromWalletID)

	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Errorf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if err := r.Debit(tx, transfer.FromWalletID, transfer.Amount); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		r.logger.Errorf("Failed to debit from_wallet: %v", err)
		return fmt.Errorf("failed to debit from_wallet: %w", err)
	}

	if err := r.Credit(tx, transfer.ToWalletID, transfer.Amount); err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		r.logger.Errorf("Failed to credit to_wallet: %v", err)
		return fmt.Errorf("failed to credit to_wallet: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *walletRepo) UpdateTransferStatus(ctx context.Context, externalRef, status string) error {
	r.logger.Infof("Updating transfer status for external_ref %s", externalRef)
	_, err := r.DB.ExecContext(ctx, "UPDATE transfers SET status = $1 WHERE external_ref = $2", status, externalRef)
	return err
}

func (r *walletRepo) Debit(txn *sqlx.Tx, walletID uuid.UUID, amount float64) error {
	r.logger.Infof("Debiting wallet %s by %f", walletID, amount)

	var balance float64
	err := txn.QueryRowContext(context.Background(), "SELECT balance FROM wallets WHERE id = ? FOR UPDATE", walletID).Scan(&balance)
	if err != nil {
		r.logger.Errorf("Failed to fetch wallet balance: %v", err)
		return fmt.Errorf("failed to fetch wallet balance: %w", err)
	}

	if balance < amount {
		r.logger.Errorf("Insufficient balance in wallet %s", walletID)
		return fmt.Errorf("insufficient balance in wallet %s", walletID)
	}

	_, err = txn.ExecContext(context.Background(), "UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, walletID)
	return err
}

func (r *walletRepo) Credit(txn *sqlx.Tx, walletID uuid.UUID, amount float64) error {
	r.logger.Infof("Crediting wallet %s by %f", walletID, amount)
	_, err := txn.ExecContext(context.Background(), "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
	return err
}
