package service

import "context"

type UserService interface {
    RegisterUser(ctx context.Context, user User) (User, error)
    GetUserByID(ctx context.Context, userID string) (User, error)
}
