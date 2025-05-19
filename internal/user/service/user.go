package service

import (
	"context"
	"github.com/markbates/goth"
	"github.com/mavrk-mose/pay/internal/user/models"
)

//go:generate mockery --name=UserService --output=./mocks --filename=user.go --with-expecter
type UserService interface {
	RegisterUser(ctx context.Context, user goth.User) (string, error)
	UpdateUser(ctx context.Context, userID string, updates models.UserUpdateRequest) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	ListUsers(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	AssignRole(ctx context.Context, userID, role string) error
	RevokeRole(ctx context.Context, userID string) error
	BanUser(ctx context.Context, userID string, reason string) error
	UnbanUser(ctx context.Context, userID string) error
}
