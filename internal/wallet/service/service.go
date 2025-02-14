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
	Transfer(req TransferRequest) error
	Withdraw(req WithdrawalRequest) error
	Deposit(req DepositRequest) error
	GetWallet(userID string) (float64, error)
	GetBalance(userID string) (float64, error)
}

type walletService struct {
	repo WalletRepo
}

func (s *walletService) CreateWallet(ctx context.Context, req CreateWalletRequest) error {
	wallet := Wallet{
		CustomerID: uuid.New(),
		Balance:    req.InitialBalance,
		Currency:   req.Currency,
	}
	return s.repo.Create(ctx, wallet)
}

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
		//TODO: implement currency conversion
		return errors.New("wallet currencies do not match")
	}

	if fromWallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	fromWallet.Balance -= req.Amount
	toWallet.Balance += req.Amount

	if err := s.repo.UpdateWalletBalance(ctx, fromWallet.ID, req.Amount); err != nil {
		return err
	}
	if err := s.repo.UpdateWalletBalance(ctx, toWallet.ID, req.Amount); err != nil {
		return err
	}
	return nil
}

func (s *walletService) Withdraw(req WithdrawalRequest) error {
	return nil
}

func (s *walletService) GetBalance(userID string) (float64, error) {
	return 10, nil

}
