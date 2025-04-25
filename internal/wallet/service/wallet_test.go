package service_test

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	models "github.com/mavrk-mose/pay/internal/wallet/models"
	"github.com/mavrk-mose/pay/internal/wallet/repository/mocks"
	"github.com/mavrk-mose/pay/internal/wallet/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWallet(t *testing.T) {
	mockRepo := new(mocks.WalletRepo)
	svc := service.NewWalletService(mockRepo)

	ctx := &gin.Context{}
	req := models.CreateWalletRequest{
		CustomerID: "user123",
		Currency:   "USD",
	}

	expectedWallet := &models.Wallet{
		UserId:    req.CustomerID,
		Currency:  req.Currency,
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("CreateWallet", ctx, mock.AnythingOfType("models.Wallet")).Return(expectedWallet, nil)

	wallet, err := svc.CreateWallet(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, req.CustomerID, wallet.UserId)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	mockRepo := new(mocks.WalletRepo)
	svc := service.NewWalletService(mockRepo)
	ctx := &gin.Context{}

	fromWallet := &models.Wallet{ID: uuid.New(), Currency: "USD", Balance: 100}
	toWallet := &models.Wallet{ID: uuid.New(), Currency: "USD", Balance: 50}

	req := models.TransferRequest{
		FromWalletID: fromWallet.ID,
		ToWalletID:   toWallet.ID,
		Amount:       30,
	}

	mockRepo.On("GetByID", ctx, fromWallet.ID.String()).Return(fromWallet, nil)
	mockRepo.On("GetByID", ctx, toWallet.ID.String()).Return(toWallet, nil)
	mockRepo.On("Debit", ctx, fromWallet.ID, float64(30)).Return(nil)
	mockRepo.On("Credit", ctx, toWallet.ID, float64(30)).Return(nil)

	err := svc.Transfer(ctx, req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	mockRepo := new(mocks.WalletRepo)
	svc := service.NewWalletService(mockRepo)
	ctx := &gin.Context{}

	fromWallet := &models.Wallet{ID: uuid.New(), Currency: "USD", Balance: 10}
	toWallet := &models.Wallet{ID: uuid.New(), Currency: "USD", Balance: 0}

	req := models.TransferRequest{
		FromWalletID: fromWallet.ID,
		ToWalletID:   toWallet.ID,
		Amount:       50,
	}

	mockRepo.On("GetByID", ctx, fromWallet.ID.String()).Return(fromWallet, nil)
	mockRepo.On("GetByID", ctx, toWallet.ID.String()).Return(toWallet, nil)

	err := svc.Transfer(ctx, req)
	assert.EqualError(t, err, "insufficient funds")
}

func TestCanWithdraw(t *testing.T) {
    mockRepo := new(mocks.WalletRepo)
    svc := service.NewWalletService(mockRepo)

    mockRepo.On("GetBalance", "user123").Return(100, nil)

    ok, err := svc.CanWithdraw("user123", 50)

    assert.NoError(t, err)
    assert.True(t, ok)
    mockRepo.AssertExpectations(t)
}


