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
	ID      string `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Message string `db:"message" json:"message"`
	Type    string `db:"type" json:"type"`
}
