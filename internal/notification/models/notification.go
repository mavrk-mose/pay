package models

import "github.com/google/uuid"

type Notification struct {
	ID      uuid.UUID `json:"id"`
	UserID  string    `json:"user_id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	Type    string    `json:"type"` // e.g., "info", "alert", "success"
}

type NotificationTemplate struct {
	ID      string `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Message string `db:"message" json:"message"`
	Type    string `db:"type" json:"type"`
}
