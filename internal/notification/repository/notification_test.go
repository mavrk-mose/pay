package repository

import (
	"context"
	"github.com/mavrk-mose/pay/internal/notification/models"
	"github.com/mavrk-mose/pay/internal/notification/repository/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewNotificationRepo(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.NotificationRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewNotificationRepo(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationRepo_GetTemplate(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx        context.Context
		templateID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.NotificationTemplate
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.NotificationRepo{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetTemplate(tt.args.ctx, tt.args.templateID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationRepo_StoreNotification(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx          context.Context
		notification models.Notification
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
			_m := &mocks.NotificationRepo{
				Mock: tt.fields.Mock,
			}
			if err := _m.StoreNotification(tt.args.ctx, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("StoreNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
