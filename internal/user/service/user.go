package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	. "github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, user goth.User) (User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx gin.Context, user goth.User) (User, error) {
	user, err := s.repo.CreateOrUpdateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
		return
	}

	token, err := GenerateJWT(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT: " + err.Error()})
		return
	}


}

func (s *userService) GetUserByID(ctx context.Context, userID string) (User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}


func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	})

	return token.SignedString(jwtSecret)
}