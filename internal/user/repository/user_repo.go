package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/mavrk-mose/pay/internal/user/models"
	"time"
)

func CreateOrUpdateUser(db *sqlx.DB, user goth.User) (*models.User, error) {
	var dbUser models.User
	query := `
		INSERT INTO users (google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (google_id) DO UPDATE
		SET name = $2, email = $3, avatar_url = $4, location = $5, language = $6, currency = $7, updated_at = $9, last_login_at = $10
		RETURNING id, google_id, name, email, avatar_url, location, language, currency, created_at, updated_at, last_login_at
	`
	now := time.Now()
	err := db.QueryRowx(
		query,
		user.UserID,
		user.Name,
		user.Email,
		user.AvatarURL,
		"",    // Default location (can be updated later)
		"sw",  // Default language (can be updated later)
		"TZS", // Default currency (can be updated later)
		now,
		now,
		now,
	).StructScan(&dbUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update user: %v", err)
	}
	return &dbUser, nil
}
