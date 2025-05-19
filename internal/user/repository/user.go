package repository

import (
	"context"
	"fmt"
	"github.com/mavrk-mose/pay/internal/api/middleware"
	userService "github.com/mavrk-mose/pay/internal/user/service"
	walletService "github.com/mavrk-mose/pay/internal/wallet/service"
	"strings"
	"time"

	wallet "github.com/mavrk-mose/pay/internal/wallet/models"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/pkg/utils"
)

type userRepo struct {
	db            *sqlx.DB
	walletService walletService.WalletService
	logger        utils.Logger
}

func NewUserService(db *sqlx.DB) userService.UserService {
	return &userRepo{
		db: db,
	}
}

func (s *userRepo) RegisterUser(ctx context.Context, user goth.User) (string, error) {
	s.logger.Infof("Registering user : %v", user)

	dbUser, err := s.CreateOrUpdateUser(ctx, user)
	if err != nil {
		s.logger.Errorf("Failed to create/update user: %v", err)
		return "", err
	}

	token, err := middleware.GenerateJWT(dbUser.ID.String())
	if err != nil {
		s.logger.Errorf("Failed to generate JWT: %v", err)
		return "", err
	}
	return token, nil
}

func (s *userRepo) CreateOrUpdateUser(ctx context.Context, user goth.User) (*models.User, error) {
	var dbUser models.User
	query := `
		INSERT INTO users (google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (google_id) DO UPDATE
		SET name = $2, email = $3, avatar_url = $4, location = $5, language = $6, currency = $7, updated_at = $9, last_login_at = $10
		RETURNING id, google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at
	`
	now := time.Now()
	err := s.db.QueryRowx(
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
		s.logger.Errorf("Failed to create/update user: %v", err)
		return nil, fmt.Errorf("failed to create/update user: %v", err)
	}

	var walletCount int
	err = s.db.Get(&walletCount, "SELECT COUNT(*) FROM wallets WHERE user_id = $1", dbUser.UserId)
	if err != nil {
		s.logger.Errorf("Failed to check wallet existence: %v", err)
		return nil, fmt.Errorf("failed to check wallet existence: %v", err)
	}

	if walletCount == 0 {
		newWallet := &wallet.Wallet{
			UserId:    dbUser.UserId,
			Balance:   0,
			Currency:  dbUser.Currency,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = s.walletService.CreateWallet(ctx, newWallet)
		if err != nil {
			s.logger.Errorf("Failed to create default newWallet: %v", err)
			return nil, fmt.Errorf("failed to create default newWallet: %v", err)
		}
	}

	return &dbUser, nil
}

func (s *userRepo) UpdateUser(ctx context.Context, userID string, updates models.UserUpdateRequest) error {
	var (
		setClauses []string
		args       []interface{}
	)

	if updates.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *updates.Name)
	}
	if updates.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, *updates.Email)
	}
	if updates.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", len(args)+1))
		args = append(args, *updates.Phone)
	}
	if updates.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *updates.IsActive)
	}

	if len(setClauses) == 0 {
		return nil // nothing to update
	}

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		len(args)+1,
	)
	args = append(args, userID)

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *userRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	query := `SELECT id, google_id, name, email, phone_number, avatar_url, location, language, currency, created_at, updated_at, last_login_at FROM users WHERE google_id = $1`
	err := s.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		s.logger.Errorf("Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}

func (s *userRepo) ListUsers(ctx context.Context, filter models.UserFilter) ([]models.User, error) {
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
	err := s.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

func (s *userRepo) AssignRole(ctx context.Context, userID, role string) error {
	s.logger.Infof("Assigning role %s to user %s", role, userID)
	_, err := s.db.ExecContext(ctx, "UPDATE users SET role = $1 WHERE id = $2", role, userID)
	return err
}

func (s *userRepo) RevokeRole(ctx context.Context, userID string) error {
	s.logger.Infof("Revoking role from user %s", userID)
	_, err := s.db.ExecContext(ctx, "UPDATE users SET role = NULL WHERE id = $1", userID)
	return err
}

func (s *userRepo) BanUser(ctx context.Context, userID, reason string) error {
	s.logger.Infof("Banning user %s due to %s", userID, reason)
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Errorf("Failed to start transaction %s", err)
		return err
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET is_active = false WHERE id = $1", userID)
	if err != nil {
		s.logger.Errorf("Failed to update user %s", err)
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO user_bans (user_id, reason) VALUES ($1, $2)", userID, reason)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func (s *userRepo) UnbanUser(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE users SET is_active = true WHERE id = $1", userID)
	return err
}
