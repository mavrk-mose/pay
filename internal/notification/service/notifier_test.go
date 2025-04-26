package service

import (
	"context"
	"github.com/mavrk-mose/pay/internal/notification/service/mocks"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewNotifier(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.Notifier
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewNotifier(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifier_Send(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx        context.Context
		user       models.User
		templateID string
		details    map[string]string
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
			_m := &mocks.Notifier{
				Mock: tt.fields.Mock,
			}
			if err := _m.Send(tt.args.ctx, tt.args.user, tt.args.templateID, tt.args.details); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
