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
	UpdateUser(ctx context.Context, userID string, updates models.UserUpdateRequest) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	ListUsers(ctx context.Context, filter models.UserFilter) ([]models.User, error)
	AssignRole(ctx context.Context, userID, role string) error
	RevokeRole(ctx context.Context, userID string) error
	BanUser(ctx context.Context, userID string, reason string) error
	UnbanUser(ctx context.Context, userID string) error
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

func (r *userRepo) UpdateUser(ctx context.Context, userID string, updates models.UserUpdateRequest) error {
	query := "UPDATE users SET "
	var args []interface{}
	argCount := 1

	if updates.Name != nil {
		query += fmt.Sprintf("name = $%d, ", argCount)
		args = append(args, *updates.Name)
		argCount++
	}
	if updates.Email != nil {
		query += fmt.Sprintf("email = $%d, ", argCount)
		args = append(args, *updates.Email)
		argCount++
	}
	if updates.Phone != nil {
		query += fmt.Sprintf("phone = $%d, ", argCount)
		args = append(args, *updates.Phone)
		argCount++
	}
	if updates.IsActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", argCount)
		args = append(args, *updates.IsActive)
		argCount++
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
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

func (r *userRepo) ListUsers(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
	query := "SELECT id, name, email, role, is_active FROM users WHERE 1=1"
	var args []interface{}
	argCount := 1

	if filter.Role != nil {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, *filter.Role)
		argCount++
	}
	if filter.Active != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filter.Active)
		argCount++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

func (r *userRepo) AssignRole(ctx context.Context, userID, role string) error {
	r.logger.Infof("Assigning role %s to user %s", role, userID)
	_, err := r.db.ExecContext(ctx, "UPDATE users SET role = $1 WHERE id = $2", role, userID)
	return err
}

func (r *userRepo) RevokeRole(ctx context.Context, userID string) error {
	r.logger.Infof("Revoking role from user %s", userID)
	_, err := r.db.ExecContext(ctx, "UPDATE users SET role = NULL WHERE id = $1", userID)
	return err
}

func (r *userRepo) BanUser(ctx context.Context, userID, reason string) error {
	r.logger.Infof("Banning user %s due to %s", userID, reason)
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Errorf("Failed to start transaction %s", err)
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET is_active = false WHERE id = $1", userID)
	if err != nil {
		r.logger.Errorf("Failed to update user %s", err)
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO user_bans (user_id, reason) VALUES ($1, $2)", userID, reason)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *userRepo) UnbanUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET is_active = true WHERE id = $1", userID)
	return err
}