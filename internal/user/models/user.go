package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserId      string    `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`                   // Full name
	Email       string    `json:"email" db:"email"`                 // Email address
	PhoneNumber string    `json:"phone_number" db:"phone_number"`   // Phone number
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`       // Profile picture URL
	Location    string    `json:"location" db:"location"`           // User's location (e.g., "New York, USA")
	Language    string    `json:"language" db:"language"`           // Preferred language (e.g., "en")
	Currency    string    `json:"currency" db:"currency"`           // Preferred currency (e.g., "USD")
	Role        string    `json:"role" db:"role"`                   // User role (e.g., "admin", "user", "moderator")
	CreatedAt   time.Time `json:"created_at" db:"created_at"`       // Timestamp when the user was created
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`       // Timestamp when the user was last updated
	LastLoginAt time.Time `json:"last_login_at" db:"last_login_at"` // Timestamp of the last login
	DeviceToken string    `json:"deviceToken" db:"device_token"`
	Onboarded   bool      `json:"onboarded" db:"onboarded"`
	Provider    string    `json:"provider" db:"provider"`
}

// role struct will have a rules of the permissions
type Permissions struct {
	Entry     int  `json:"entry" db:"entry"`
	AddFlag   bool `json:"add_flag" db:"add_flag"`
	AdminFlag bool `json:"admin_flag" db:"admin_flag"`
}
