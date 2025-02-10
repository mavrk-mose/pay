package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mavrk-mose/pay/internal/model"
	"github.com/mavrk-mose/pay/internal/repository"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req model.CreateWalletRequest) (model.Wallet, error)
	Transfer(req model.TransferRequest) error
	Withdraw(req model.WithdrawalRequest) error
	Deposit(req model.DepositRequest) error
	GetWallet(userID string) (float64, error)
}

type walletService struct {
	repo repository.WalletRepo
}

func NewWalletService(repo repository.WalletRepo) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateWallet(ctx context.Context, req model.CreateWalletRequest) error {
	wallet := model.Wallet{
		CustomerID: uuid.New(),
		Balance:    req.InitialBalance,
		Currency:   req.Currency,
	}
	return s.repo.Create(ctx, wallet)
}

func (s *walletService) Transfer(ctx context.Context, req model.TransferRequest) error {
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

	if err := s.repo.UpdateWalletBalance(fromWallet); err != nil {
		return err
	}
	if err := s.repo.UpdateWalletBalance(toWallet); err != nil {
		return err
	}
	return nil
}

func (s *walletService) Withdraw(req model.WithdrawalRequest) error {

}
