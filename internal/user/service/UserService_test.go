package service

import (
	"context"
	"github.com/markbates/goth"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/service/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		t interface {
			mock.TestingT
			Cleanup(func())
		}
	}
	tests := []struct {
		name string
		args args
		want *mocks.UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mocks.NewUserService(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_AssignRole(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		userID string
		role   string
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			if err := _m.AssignRole(tt.args.ctx, tt.args.userID, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("AssignRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_BanUser(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		userID string
		reason string
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			if err := _m.BanUser(tt.args.ctx, tt.args.userID, tt.args.reason); (err != nil) != tt.wantErr {
				t.Errorf("BanUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.GetUserByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		filter models.UserFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.ListUsers(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_RegisterUser(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx  context.Context
		user goth.User
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			got, err := _m.RegisterUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_RevokeRole(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		userID string
		role   string
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			if err := _m.RevokeRole(tt.args.ctx, tt.args.userID, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("RevokeRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UnbanUser(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx    context.Context
		userID string
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			if err := _m.UnbanUser(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("UnbanUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type fields struct {
		Mock mock.Mock
	}
	type args struct {
		ctx     context.Context
		userID  string
		updates models.UserUpdateRequest
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
			_m := &mocks.UserService{
				Mock: tt.fields.Mock,
			}
			if err := _m.UpdateUser(tt.args.ctx, tt.args.userID, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
