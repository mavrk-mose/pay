package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
	. "github.com/mavrk-mose/pay/internal/wallet/repository"
)

// Wallet module (tracks balances)

type WalletService interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest) (Wallet, error)
	Transfer(ctx context.Context, req TransferRequest) error
	Withdraw(ctx context.Context, req WithdrawalRequest) error
	Deposit(ctx context.Context, req DepositRequest) error
	GetWallet(ctx context.Context, userID string) (Wallet, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error)
}

type walletService struct {
	repo WalletRepo
}

func NewWalletService(repo WalletRepo) WalletService {
	return &walletService{repo: repo}
}

// CreateWallet creates a new wallet for a user
func (s *walletService) CreateWallet(ctx context.Context, req CreateWalletRequest) (Wallet, error) {
	wallet := Wallet{
		CustomerID: uuid.New(),
		Balance:    req.InitialBalance,
		Currency:   req.Currency,
	}
	err := s.repo.Create(ctx, wallet)
	return wallet, err
}

// Transfer moves funds from one wallet to another
func (s *walletService) Transfer(ctx context.Context, req TransferRequest) error {
	fromWallet, err := s.repo.GetByID(ctx, req.FromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := s.repo.GetByID(ctx, req.ToWalletID)
	if err != nil {
		return err
	}

	if fromWallet.Currency != toWallet.Currency {
		return errors.New("wallet currencies do not match")
	}

	if fromWallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	fromWallet.Balance -= req.Amount
	toWallet.Balance += req.Amount

	if err := s.repo.UpdateWalletBalance(ctx, fromWallet.ID, -req.Amount); err != nil {
		return err
	}
	if err := s.repo.UpdateWalletBalance(ctx, toWallet.ID, req.Amount); err != nil {
		return err
	}
	return nil
}

// Withdraw handles withdrawing funds from a wallet
func (s *walletService) Withdraw(ctx context.Context, req WithdrawalRequest) error {
	wallet, err := s.repo.GetByID(ctx, req.WalletID)
	if err != nil {
		return err
	}

	if wallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	wallet.Balance -= req.Amount
	return s.repo.UpdateWalletBalance(ctx, wallet.ID, -req.Amount)
}

// Deposit adds funds to a wallet
func (s *walletService) Deposit(ctx context.Context, req DepositRequest) error {
	wallet, err := s.repo.GetByID(ctx, req.WalletID)
	if err != nil {
		return err
	}

	wallet.Balance += req.Amount
	return s.repo.UpdateWalletBalance(ctx, wallet.ID, req.Amount)
}

// GetWallet retrieves a wallet by user ID
func (s *walletService) GetWallet(ctx context.Context, userID string) (Wallet, error) {
	return s.repo.GetByID(ctx, userID)
}

// GetBalance returns the balance of a wallet
func (s *walletService) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// UpdateBalance updates the balance of a wallet
func (s *walletService) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	return s.repo.UpdateWalletBalance(ctx, walletID, amount)
}
