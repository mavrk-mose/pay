package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/mavrk-mose/pay/internal/wallet/models"
	. "github.com/mavrk-mose/pay/internal/wallet/repository"
	"time"
)

//go:generate mockery --name=WalletService --output=./mocks --filename=wallet.go --with-expecter
type WalletService interface {
	CreateWallet(ctx *gin.Context, req CreateWalletRequest) (Wallet, error)
	Transfer(ctx *gin.Context, req TransferRequest) error
	DebitWallet(ctx *gin.Context, req WalletTransactionRequest) error
	CreditWallet(ctx *gin.Context, req WalletTransactionRequest) error
	GetWallet(ctx *gin.Context, userID string) (Wallet, error)
	UpdateBalance(ctx *gin.Context, walletID uuid.UUID, amount float64) error
	GetBalance(ctx *gin.Context, walletID uuid.UUID) (float64, error)
	GetUserWallets(ctx *gin.Context, id string) ([]Wallet, error)
	DeleteWallet(c *gin.Context, walletID string) any
	CanWithdraw(request string, i int) (any, error)
}

type walletService struct {
	repo WalletRepo
}

func NewWalletService(repo WalletRepo) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) CanWithdraw(request string, i int) (any, error) {
	panic("unimplemented")
}

// CreateWallet creates a new wallet for a user
func (s *walletService) CreateWallet(ctx *gin.Context, req CreateWalletRequest) (Wallet, error) {
	newWallet := Wallet{
		UserId:    req.CustomerID,
		Balance:   0.00,
		Currency:  req.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	wallet, err := s.repo.CreateWallet(ctx, newWallet)
	if err != nil {
		return Wallet{}, err
	}

	return *wallet, nil
}

// Transfer moves funds from one wallet to another
func (s *walletService) Transfer(ctx *gin.Context, req TransferRequest) error {
	fromWallet, err := s.repo.GetByID(ctx, req.FromWalletID.String())
	if err != nil {
		return err
	}

	toWallet, err := s.repo.GetByID(ctx, req.ToWalletID.String())
	if err != nil {
		return err
	}

	if fromWallet.Currency != toWallet.Currency {
		return errors.New("wallet currencies do not match")
	}

	if fromWallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	if err := s.repo.Debit(ctx, req.FromWalletID, req.Amount); err != nil {
		return err
	}
	if err := s.repo.Credit(ctx, req.ToWalletID, req.Amount); err != nil {
		return err
	}
	return nil
}

// Withdraw handles withdrawing funds from a wallet
func (s *walletService) DebitWallet(ctx *gin.Context, req WalletTransactionRequest) error {
	wallet, err := s.repo.GetByID(ctx, req.WalletID)
	if err != nil {
		return err
	}

	if wallet.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	return s.repo.Debit(ctx, uuid.MustParse(req.WalletID), req.Amount)
}

// Deposit adds funds to a wallet
func (s *walletService) CreditWallet(ctx *gin.Context, req WalletTransactionRequest) error {
	return s.repo.Credit(ctx, uuid.MustParse(req.WalletID), req.Amount)
}

// GetWallet retrieves a wallet by user ID
func (s *walletService) GetWallet(ctx *gin.Context, userID string) (Wallet, error) {
	return s.repo.GetByID(ctx, userID)
}

// GetBalance returns the balance of a wallet
func (s *walletService) GetBalance(ctx *gin.Context, walletID uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetByID(ctx, walletID.String())
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

// UpdateBalance updates the balance of a wallet
func (s *walletService) UpdateBalance(ctx *gin.Context, walletID uuid.UUID, amount float64) error {
	return s.repo.Credit(ctx, walletID, amount)
}

func (s *walletService) GetUserWallets(ctx *gin.Context, id string) ([]Wallet, error) {
	return s.repo.GetUserWallets(ctx.Request.Context(), id)
}

func (s *walletService) DeleteWallet(c *gin.Context, walletID string) any {
	panic("unimplemented")
}
