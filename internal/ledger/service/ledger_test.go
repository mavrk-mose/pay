package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/mavrk-mose/pay/internal/ledger/models"
	"github.com/mavrk-mose/pay/internal/ledger/service/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestLedgerService_CreateTransactionWithEntries(t *testing.T) {
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
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			if err := _m.CreateTransactionWithEntries(tt.args.ctx, tt.args.txn, tt.args.entries); (err != nil) != tt.wantErr {
				t.Errorf("CreateTransactionWithEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLedgerService_CreateTransactionWithEntries_Call_Return(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		_a0 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_CreateTransactionWithEntries_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_CreateTransactionWithEntries_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Return(tt.args._a0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_CreateTransactionWithEntries_Call_Run(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(ctx *gin.Context, txn *sqlx.Tx, entries []models.Transaction)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_CreateTransactionWithEntries_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_CreateTransactionWithEntries_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Run(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_CreateTransactionWithEntries_Call_RunAndReturn(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(*gin.Context, *sqlx.Tx, []models.Transaction) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_CreateTransactionWithEntries_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_CreateTransactionWithEntries_Call{
				Call: tt.fields.Call,
			}
			if got := _c.RunAndReturn(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAndReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_EXPECT(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	tests := []struct {
		name   string
		fields fields
		want   *mocks.LedgerService_Expecter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			if got := _m.EXPECT(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EXPECT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_Expecter_CreateTransactionWithEntries(t *testing.T) {
	type fields struct {
		mock *mock.Mock
	}
	type args struct {
		ctx     interface{}
		txn     interface{}
		entries interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_CreateTransactionWithEntries_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_e := &mocks.LedgerService_Expecter{
				mock: tt.fields.mock,
			}
			if got := _e.CreateTransactionWithEntries(tt.args.ctx, tt.args.txn, tt.args.entries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTransactionWithEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_Expecter_FetchTransactionsWithChecksum(t *testing.T) {
	type fields struct {
		mock *mock.Mock
	}
	type args struct {
		db       interface{}
		date     interface{}
		provider interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_FetchTransactionsWithChecksum_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_e := &mocks.LedgerService_Expecter{
				mock: tt.fields.mock,
			}
			if got := _e.FetchTransactionsWithChecksum(tt.args.db, tt.args.date, tt.args.provider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchTransactionsWithChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_Expecter_GetTransactionByID(t *testing.T) {
	type fields struct {
		mock *mock.Mock
	}
	type args struct {
		transactionID interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_GetTransactionByID_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_e := &mocks.LedgerService_Expecter{
				mock: tt.fields.mock,
			}
			if got := _e.GetTransactionByID(tt.args.transactionID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTransactionByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_Expecter_RecordTransaction(t *testing.T) {
	type fields struct {
		mock *mock.Mock
	}
	type args struct {
		ctx interface{}
		txn interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_RecordTransaction_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_e := &mocks.LedgerService_Expecter{
				mock: tt.fields.mock,
			}
			if got := _e.RecordTransaction(tt.args.ctx, tt.args.txn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RecordTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_Expecter_UpdateTransactionStatus(t *testing.T) {
	type fields struct {
		mock *mock.Mock
	}
	type args struct {
		ctx         interface{}
		externalRef interface{}
		status      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_UpdateTransactionStatus_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_e := &mocks.LedgerService_Expecter{
				mock: tt.fields.mock,
			}
			if got := _e.UpdateTransactionStatus(tt.args.ctx, tt.args.externalRef, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateTransactionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_FetchTransactionsWithChecksum(t *testing.T) {
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
			_m := &mocks.LedgerService{
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

func TestLedgerService_FetchTransactionsWithChecksum_Call_Return(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		_a0 map[string]string
		_a1 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_FetchTransactionsWithChecksum_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_FetchTransactionsWithChecksum_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Return(tt.args._a0, tt.args._a1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_FetchTransactionsWithChecksum_Call_Run(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(db *sqlx.DB, date string, provider string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_FetchTransactionsWithChecksum_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_FetchTransactionsWithChecksum_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Run(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_FetchTransactionsWithChecksum_Call_RunAndReturn(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(*sqlx.DB, string, string) (map[string]string, error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_FetchTransactionsWithChecksum_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_FetchTransactionsWithChecksum_Call{
				Call: tt.fields.Call,
			}
			if got := _c.RunAndReturn(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAndReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestLedgerService_GetTransactionByID_Call_Return(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		_a0 models.Transaction
		_a1 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_GetTransactionByID_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_GetTransactionByID_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Return(tt.args._a0, tt.args._a1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_GetTransactionByID_Call_Run(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(transactionID string)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_GetTransactionByID_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_GetTransactionByID_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Run(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_GetTransactionByID_Call_RunAndReturn(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(string) (models.Transaction, error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_GetTransactionByID_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_GetTransactionByID_Call{
				Call: tt.fields.Call,
			}
			if got := _c.RunAndReturn(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAndReturn() = %v, want %v", got, tt.want)
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
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.RecordTransaction(tt.args.ctx, tt.args.txn)
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

func TestLedgerService_RecordTransaction_Call_Return(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		_a0 string
		_a1 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_RecordTransaction_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_RecordTransaction_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Return(tt.args._a0, tt.args._a1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_RecordTransaction_Call_Run(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(ctx *gin.Context, txn models.Transaction)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_RecordTransaction_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_RecordTransaction_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Run(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_RecordTransaction_Call_RunAndReturn(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(*gin.Context, models.Transaction) (string, error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_RecordTransaction_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_RecordTransaction_Call{
				Call: tt.fields.Call,
			}
			if got := _c.RunAndReturn(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAndReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_UpdateTransactionStatus(t *testing.T) {
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
			_m := &mocks.LedgerService{
				Mock: tt.fields.Mock,
			}
			if err := _m.UpdateTransactionStatus(tt.args.ctx, tt.args.externalRef, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLedgerService_UpdateTransactionStatus_Call_Return(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		_a0 error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_UpdateTransactionStatus_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_UpdateTransactionStatus_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Return(tt.args._a0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Return() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_UpdateTransactionStatus_Call_Run(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(ctx *gin.Context, externalRef string, status models.TransactionStatus)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_UpdateTransactionStatus_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_UpdateTransactionStatus_Call{
				Call: tt.fields.Call,
			}
			if got := _c.Run(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLedgerService_UpdateTransactionStatus_Call_RunAndReturn(t *testing.T) {
	type fields struct {
		Call *mock.Call
	}
	type args struct {
		run func(*gin.Context, string, models.TransactionStatus) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mocks.LedgerService_UpdateTransactionStatus_Call
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := &mocks.LedgerService_UpdateTransactionStatus_Call{
				Call: tt.fields.Call,
			}
			if got := _c.RunAndReturn(tt.args.run); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunAndReturn() = %v, want %v", got, tt.want)
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
