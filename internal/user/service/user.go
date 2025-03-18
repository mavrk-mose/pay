package service

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	. "github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

var (
	jwtSecret      []byte
	expirationTime time.Time
)

type UserService interface {
	RegisterUser(ctx context.Context, user goth.User) (token string, err error)
	GetUserByID(ctx context.Context, userID string) (User, error)
}

type userService struct {
	repo repository.UserRepository
	logger utils.Logger
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, user goth.User) (string, error) {
	dbUser, err := s.repo.CreateOrUpdateUser(ctx, user)
	if err != nil {
		s.logger.Errorf("Failed to create/update user: %v", err)
		return "", err
	}

	token, err := GenerateJWT(dbUser.ID.String())
	if err != nil {
		s.logger.Errorf("Failed to generate JWT: %v", err)
		return "", err
	}
	return token, nil

}

func (s *userService) GetUserByID(ctx context.Context, userID string) (User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}


func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	})

	return token.SignedString(jwtSecret)
}

