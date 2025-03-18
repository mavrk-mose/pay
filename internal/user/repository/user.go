package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	models "github.com/mavrk-mose/pay/internal/user/models"
	walletRepo "github.com/mavrk-mose/pay/internal/wallet/repository"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type UserRepository interface {
	CreateOrUpdateUser(ctx context.Context, user goth.User) (*models.User, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
}

type userRepo struct {
	db         *sqlx.DB
	walletRepo walletRepo.WalletRepo
	logger     utils.Logger
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateOrUpdateUser(ctx context.Context, user goth.User) (*models.User, error) {
	var dbUser models.User
	query := `
		INSERT INTO users (google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (google_id) DO UPDATE
		SET name = $2, email = $3, avatar_url = $4, location = $5, language = $6, currency = $7, updated_at = $9, last_login_at = $10
		RETURNING id, google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at
	`
	now := time.Now()
	err := r.db.QueryRowx(
		query,
		user.UserID,
		user.Name,
		user.Email,
		user.AvatarURL,
		"",
		"sw",
		"TZS",
		now,
		now,
		now,
	).StructScan(&dbUser)
	if err != nil {
		r.logger.Errorf("Failed to create/update user: %v", err)
		return nil, fmt.Errorf("failed to create/update user: %v", err)
	}
	return &dbUser, nil

	var walletCount int
	err = r.db.Get(&walletCount, "SELECT COUNT(*) FROM wallets WHERE user_id = $1", dbUser.UserId)
	if err != nil {
		r.logger.Errorf("Failed to check wallet existence: %v", err)
		return nil, fmt.Errorf("failed to check wallet existence: %v", err)
	}

	if walletCount == 0 {
		_, err = r.walletRepo.CreateWallet(ctx, dbUser.UserId, dbUser.Currency)
		if err != nil {
			r.logger.Errorf("Failed to create default wallet: %v", err)
			return nil, fmt.Errorf("failed to create default wallet: %v", err)
		}
	}

	return &dbUser, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, google_id, name, email, phone_number, avatar_url, location, language, currency, created_at, updated_at, last_login_at FROM users WHERE google_id = $1`
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		r.logger.Errorf("Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}
