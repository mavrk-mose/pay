package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	wallet "github.com/mavrk-mose/pay/internal/wallet/models"
	"github.com/mavrk-mose/pay/internal/wallet/service/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewWalletService(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.WalletService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewWalletService(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWalletService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_CanWithdraw(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		request string
		i       int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.CanWithdraw(tt.args.request, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("CanWithdraw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CanWithdraw() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_CreateWallet(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		req wallet.CreateWalletRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    wallet.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.CreateWallet(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_CreditWallet(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		req wallet.WalletTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			if err := _m.CreditWallet(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CreditWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_DebitWallet(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		req wallet.WalletTransactionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			if err := _m.DebitWallet(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("DebitWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_DeleteWallet(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		c        *gin.Context
		walletID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			if got := _m.DeleteWallet(tt.args.c, tt.args.walletID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteWallet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetBalance(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx      *gin.Context
		walletID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetBalance(tt.args.ctx, tt.args.walletID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetUserWallets(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []wallet.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetUserWallets(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserWallets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserWallets() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_GetWallet(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    *gin.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    wallet.Wallet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetWallet(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWallet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalletService_Transfer(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		req wallet.TransferRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			if err := _m.Transfer(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_UpdateBalance(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx      *gin.Context
		walletID uuid.UUID
		amount   float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.WalletService{
				Mock: tt.fields.Mock,
			}
			if err := _m.UpdateBalance(tt.args.ctx, tt.args.walletID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
