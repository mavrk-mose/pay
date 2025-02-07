package repository

import (
    "context"
    "database/sql"
    "errors"
    "github.com/jmoiron/sqlx"
    "github.com/mavrk-mose/pay/internal/model"
)

type WalletRepo struct {
    DB *sqlx.DB
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
