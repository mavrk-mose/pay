package models

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Channel   string    `json:"channel" db:"channel"` // sms, email, push, web
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"` // info, alert, success, etc.
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type NotificationTemplate struct {
	Id        string   `json:"id" db:"id"`
	Title     string   `json:"title" db:"title"`
	Subject   string   `json:"subject" db:"subject"`
	Message   string   `json:"message" db:"message"`
	Type      string   `json:"type" db:"type"`
	Channel   string   `json:"channel" db:"channel"`
	Variables []string `json:"variables" db:"variables"`
	Metadata  struct {
		Language string `json:"language" db:"language"`
	} `json:"metadata" db:"metadata"`
}
