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

//go:generate mockery --name=UserService --output=./mocks --filename=user.go --with-expecter
type UserService interface {
	RegisterUser(ctx context.Context, user goth.User) (token string, err error)
	UpdateUser(ctx context.Context, userID string, updates UserUpdateRequest) error
	GetUserByID(ctx context.Context, userID string) (User, error)
	ListUsers(ctx context.Context, filter UserFilter) ([]User, error)
	AssignRole(ctx context.Context, userID, role string) error
	RevokeRole(ctx context.Context, userID, role string) error
	BanUser(ctx context.Context, userID string, reason string) error
	UnbanUser(ctx context.Context, userID string) error
}

type userService struct {
	repo   repository.UserRepository
	logger utils.Logger
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, user goth.User) (string, error) {
	s.logger.Infof("Registering user : %v", user)

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

func (s *userService) UpdateUser(ctx context.Context, userID string, updates UserUpdateRequest) error {
	s.logger.Infof("Updating user : %v", userID)

	return s.repo.UpdateUser(ctx, userID, updates)
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (User, error) {
	s.logger.Infof("Fetching user : %v", userID)

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}
	return *user, nil
}

func (s *userService) ListUsers(ctx context.Context, filter UserFilter) ([]User, error) {
	s.logger.Infof("Listing users")

	return s.repo.ListUsers(ctx, filter)
}

func (s *userService) AssignRole(ctx context.Context, userID, role string) error {
	s.logger.Infof("Assigning role %s to user %s", role, userID)

	return s.repo.AssignRole(ctx, userID, role)
}

func (s *userService) RevokeRole(ctx context.Context, userID, role string) error {
	s.logger.Infof("Revoking role %s from user %s", role, userID)

	return s.repo.RevokeRole(ctx, userID)
}

func (s *userService) BanUser(ctx context.Context, userID string, reason string) error {
	s.logger.Infof("Banning user %s", userID)

	return s.repo.BanUser(ctx, userID, reason)
}

func (s *userService) UnbanUser(ctx context.Context, userID string) error {
	s.logger.Infof("Unbanning user %s", userID)

	return s.repo.UnbanUser(ctx, userID)
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	})

	return token.SignedString(jwtSecret)
}
