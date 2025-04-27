package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/service/mocks"
	"github.com/stretchr/testify/mock"
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
		{
            name: "create new user service",
            args: args{
                t: &testing.T{},
            },
            want: &mocks.UserService{},
        },
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
		{
            name: "successful role assignment",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("AssignRole", mock.Anything, "user123", "admin").Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                role:   "admin",
            },
            wantErr: false,
        },
        {
            name: "failed role assignment",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("AssignRole", mock.Anything, "user123", "admin").Return(errors.New("assignment failed"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                role:   "admin",
            },
            wantErr: true,
        },
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
		{
            name: "successful user ban",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("BanUser", mock.Anything, "user123", "violation").Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                reason: "violation",
            },
            wantErr: false,
        },
        {
            name: "failed user ban",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("BanUser", mock.Anything, "user123", "violation").Return(errors.New("ban failed"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                reason: "violation",
            },
            wantErr: true,
        },
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
		{
            name: "successful user retrieval",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("GetUserByID", mock.Anything, "user123").Return(models.User{ID: uuid.New(), Name: "John Doe"}, nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
            },
            want:    models.User{ID: uuid.New(), Name: "John Doe"},
            wantErr: false,
        },
        {
            name: "failed user retrieval",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("GetUserByID", mock.Anything, "user123").Return(models.User{}, errors.New("user not found"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
            },
            want:    models.User{},
            wantErr: true,
        },
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
		{
            name: "successful user list retrieval",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("ListUsers", mock.Anything, mock.Anything).Return([]models.User{
                        {ID: uuid.New(), Name: "John Doe"},
                        {ID: uuid.New(), Name: "Jane Doe"},
                    }, nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                filter: models.UserFilter{},
            },
            want: []models.User{
                {ID: uuid.New(), Name: "John Doe"},
                {ID: uuid.New(), Name: "Jane Doe"},
            },
            wantErr: false,
        },
        {
            name: "failed user list retrieval",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("ListUsers", mock.Anything, mock.Anything).Return(nil, errors.New("list failed"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                filter: models.UserFilter{},
            },
            want:    nil,
            wantErr: true,
        },
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
		{
            name: "successful user registration",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("RegisterUser", mock.Anything, mock.Anything).Return("user123", nil)
                    return m
                }(),
            },
            args: args{
                ctx:  context.TODO(),
                user: goth.User{UserID: "user123", Name: "John Doe"},
            },
            want:    "user123",
            wantErr: false,
        },
        {
            name: "failed user registration",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("RegisterUser", mock.Anything, mock.Anything).Return("", errors.New("registration failed"))
                    return m
                }(),
            },
            args: args{
                ctx:  context.TODO(),
                user: goth.User{UserID: "user123", Name: "John Doe"},
            },
            want:    "",
            wantErr: true,
        },
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
		{
            name: "successful role revocation",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("RevokeRole", mock.Anything, "user123", "admin").Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                role:   "admin",
            },
            wantErr: false,
        },
        {
            name: "failed role revocation",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("RevokeRole", mock.Anything, "user123", "admin").Return(errors.New("revocation failed"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
                role:   "admin",
            },
            wantErr: true,
        },
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
		{
            name: "successful user unban",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("UnbanUser", mock.Anything, "user123").Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
            },
            wantErr: false,
        },
        {
            name: "failed user unban",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("UnbanUser", mock.Anything, "user123").Return(errors.New("unban failed"))
                    return m
                }(),
            },
            args: args{
                ctx:    context.TODO(),
                userID: "user123",
            },
            wantErr: true,
        },
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
		{
            name: "successful user update",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("UpdateUser", mock.Anything, "user123", mock.Anything).Return(nil)
                    return m
                }(),
            },
            args: args{
                ctx:     context.TODO(),
                userID:  "user123",
                updates: models.UserUpdateRequest{Name: "Updated Name"},
            },
            wantErr: false,
        },
        {
            name: "failed user update",
            fields: fields{
                Mock: func() mock.Mock {
                    m := mock.Mock{}
                    m.On("UpdateUser", mock.Anything, "user123", mock.Anything).Return(errors.New("update failed"))
                    return m
                }(),
            },
            args: args{
                ctx:     context.TODO(),
                userID:  "user123",
                updates: models.UserUpdateRequest{Name: "Updated Name"},
            },
            wantErr: true,
        },
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
