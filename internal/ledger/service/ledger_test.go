package service

import (
	"github.com/gin-gonic/gin"
	"github.com/mavrk-mose/pay/internal/ledger/models"
	"github.com/mavrk-mose/pay/internal/ledger/service/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestLedgerService_GetTransactionByID(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		transactionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetTransactionByID(tt.args.transactionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactionByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_RecordTransaction(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx *gin.Context
		txn models.Transaction
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
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			if err := _m.RecordTransaction(tt.args.ctx, tt.args.txn); (err != nil) != tt.wantErr {
				t.Errorf("RecordTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLedgerService(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.LedgerService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewLedgerService(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLedgerService() = %v, want %v", got, tt.want)
			}
		})
	}
}
