package handler

import (
	"fmt"
	"github.com/markbates/goth/gothic"
	"github.com/mavrk-mose/pay/internal/user/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

var (
	oauthStateString = "}4PYRBlq{~m7)@wt%7jHfjo]8QyHaL6QxkwoB" // Change this to a secure random string in production
)

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	goth.UseProviders(
		google.New(clientID, clientSecret, redirectURL, "email", "profile"),
	)
}

func HandleGoogleLogin(c *gin.Context) {
	// Start the OAuth process with Goth
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func HandleGoogleCallback(c *gin.Context, db *sqlx.DB) {
	// Complete the OAuth process and get the user
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete user auth"})
		return
	}

	// Create or update the user in the database
	dbUser, err := CreateOrUpdateUser(db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/update user"})
		return
	}

	// Set user session or token (e.g., JWT)
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": dbUser})
}

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
