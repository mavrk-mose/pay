package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/mavrk-mose/pay/internal/notification/service/mocks"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/stretchr/testify/mock"
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
		{
            name: "create new notifier",
            args: args{
                t: &testing.T{},
            },
            want: &mocks.Notifier{},
        },
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
		{
            name: "successful notification send",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:        context.TODO(),
                user:       models.User{ID: uuid.New(), Name: "John Doe"},
                templateID: "welcome_email",
                details:    map[string]string{"subject": "Welcome", "body": "Hello, John!"},
            },
            wantErr: false,
        },
        {
            name: "failed notification send",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("send failed"))
                    return m
                }(),
            },
            args: args{
                ctx:        context.TODO(),
                user:       models.User{ID: uuid.New(), Name: "Jane Doe"},
                templateID: "error_email",
                details:    map[string]string{"subject": "Error", "body": "Something went wrong."},
            },
            wantErr: true,
        },
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
