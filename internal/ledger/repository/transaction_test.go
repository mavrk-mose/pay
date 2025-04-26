package ledger

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/internal/ledger/models"
	"github.com/mavrk-mose/pay/internal/ledger/repository/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewTransactionRepo(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.TransactionRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewTransactionRepo(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepo_CreateTransactionWithEntries(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx     *gin.Context
		txn     *sqlx.Tx
		entries []models.Transaction
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
			_m := &mocks.TransactionRepo{
				Mock: tt.fields.Mock,
			}
			if err := _m.CreateTransactionWithEntries(tt.args.ctx, tt.args.txn, tt.args.entries); (err != nil) != tt.wantErr {
				t.Errorf("CreateTransactionWithEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionRepo_FetchTransactionsWithChecksum(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		db       *sqlx.DB
		date     string
		provider string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.TransactionRepo{
				Mock: tt.fields.Mock,
			}
			got, err := _m.FetchTransactionsWithChecksum(tt.args.db, tt.args.date, tt.args.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionsWithChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchTransactionsWithChecksum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepo_RecordTransaction(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx           *gin.Context
		payerWalletID int64
		payeeWalletID int64
		amount        float64
		currency      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.TransactionRepo{
				Mock: tt.fields.Mock,
			}
			got, err := _m.RecordTransaction(tt.args.ctx, tt.args.payerWalletID, tt.args.payeeWalletID, tt.args.amount, tt.args.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RecordTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionRepo_UpdateTransactionStatus(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx         *gin.Context
		externalRef string
		status      models.TransactionStatus
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
			_m := &mocks.TransactionRepo{
				Mock: tt.fields.Mock,
			}
			if err := _m.UpdateTransactionStatus(tt.args.ctx, tt.args.externalRef, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
