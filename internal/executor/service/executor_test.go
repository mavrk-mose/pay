package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mavrk-mose/pay/internal/executor/service/mocks"
	"github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/stretchr/testify/mock"
)

func TestExecutorService_ExecutePayment(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		order models.PaymentIntent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "successful payment execution",
			fields: fields{
				Mock: func() mock.Mock {
					m := mock.Mock{}
					m.On("ExecutePayment", mock.Anything).Return("success", nil)
					return m
				}(),
			},
			args: args{
				order: models.PaymentIntent{ID: "123", Amount: 100},
			},
			want:    "success",
			wantErr: false,
		},
		{
			name: "failed payment execution",
			fields: fields{
				Mock: func() mock.Mock {
					m := mock.Mock{}
					m.On("ExecutePayment", mock.Anything).Return(nil, errors.New("execution failed"))
					return m
				}(),
			},
			args: args{
				order: models.PaymentIntent{ID: "124", Amount: 200},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.ExecutorService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.ExecutePayment(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecutePayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecutePayment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutorService_RecordPaymentOrder(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		order models.PaymentIntent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful record payment order",
			fields: fields{
				Mock: func() mock.Mock {
					m := mock.Mock{}
					m.On("RecordPaymentOrder", mock.Anything).Return(nil)
					return m
				}(),
			},
			args: args{
				order: models.PaymentIntent{ID: "125", Amount: 300},
			},
			wantErr: false,
		},
		{
			name: "failed record payment order",
			fields: fields{
				Mock: func() mock.Mock {
					m := mock.Mock{}
					m.On("RecordPaymentOrder", mock.Anything).Return(errors.New("record failed"))
					return m
				}(),
			},
			args: args{
				order: models.PaymentIntent{ID: "126", Amount: 400},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.ExecutorService{
				Mock: tt.fields.Mock,
			}
			if err := _m.RecordPaymentOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("RecordPaymentOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPaymentExecutorService(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.ExecutorService
	}{
		{
			name: "create new payment executor service",
			args: args{
				t: &testing.T{},
			},
			want: &mocks.ExecutorService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mocks.NewPaymentExecutorService(tt.args.t)
			if got == nil {
				t.Errorf("NewPaymentExecutorService() returned nil")
			}
		})
	}
}
