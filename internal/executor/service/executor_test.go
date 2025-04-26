package service

import (
	"github.com/mavrk-mose/pay/internal/executor/service/mocks"
	"github.com/mavrk-mose/pay/internal/payment/models"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewPaymentExecutorService(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaymentExecutorService() = %v, want %v", got, tt.want)
			}
		})
	}
}
