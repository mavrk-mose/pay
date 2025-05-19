package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
)

//go:generate mockery --name=WalletService --output=./mocks --filename=wallet.go --with-expecter
type WalletService interface {
	CreateWallet(ctx context.Context, wallet *Wallet) (*Wallet, error)
	DeleteWallet(c context.Context, walletID string) error
	GetUserWallets(ctx context.Context, userID string) ([]Wallet, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)
	GetByID(ctx context.Context, walletID uuid.UUID) (Wallet, error)
	CreateTransfer(ctx context.Context, transfer *TransferRequest) error
	UpdateTransferStatus(ctx context.Context, externalRef, status string) error
	Debit(txn *sqlx.Tx, walletID uuid.UUID, amount float64) error
	Credit(txn *sqlx.Tx, walletID uuid.UUID, amount float64) error
}
