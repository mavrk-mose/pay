package services

import (
	"errors"
	"github.com/mavrk-mose/pay/internal/model"
	"github.com/mavrk-mose/pay/internal/repository"
)

type WalletService interface {
	CreateWallet(req model.CreateWalletRequest) (model.Wallet, error)
	Transfer(req model.TransferRequest) error
	Withdraw(req model.WithdrawalRequest) error
	Deposit(req model.DepositRequest) error
	GetWallet(userID string) (float64, error)
}

type walletService struct {
	repo repository.WalletRepo
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CreateWallet(req model.CreateWalletRequest) (model.Wallet, error) {
	wallet := model.Wallet{
		CustomerID: req.CustomerID,
		Balance:    req.InitialBalance,
		Currency:   req.Currency,
	}
	return s.repo.Create(wallet)
}

func (s *walletService) Transfer(req model.TransferRequest) error {
	fromWallet, err := s.repo.GetByID(req.FromWalletID)
	if err != nil {
		return err
	}

	toWallet, err := s.repo.GetByID(req.ToWalletID)
	if err != nil {
		return err
	}

	// Optional: Ensure wallets use the same currency.
	if fromWallet.Currency != toWallet.Currency {
		//TODO: implement currency conversion
		return errors.New("wallet currencies do not match")
	}

	if fromWallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	fromWallet.Balance -= req.Amount
	toWallet.Balance += req.Amount

	if err := s.repo.Update(fromWallet); err != nil {
		return err
	}
	if err := s.repo.Update(toWallet); err != nil {
		return err
	}
	return nil
}

func (s *walletService) Withdraw(req model.WithdrawalRequest) error {

}
